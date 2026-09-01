package lexer

import (
	"fmt"
	"minimal/minimal-core/built-in/messenger"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"testing"
	"unsafe"
)

// TODO use 32-bit integers

const (
    UNKNOWN TokenType = iota
    // Always keep END at the end
    END
)

type TokenType uint32

type Token struct {
    Type  TokenType
    Value string
}

type TokenTypeMetadata struct {
    NounPhrase  string
    DebugName   string
}

type LexerScheme struct {
    maxSafePeek   uint32
    matchers      []Matcher
    lastTokenType TokenType
    metadata      []TokenTypeMetadata
}

type Lexer struct {
    Matchers    []Matcher
    Data        string
    Position    uint32
    buffer      []Token
    read        uint32
    write       uint32
    spillage    []Token
    maxSafePeek uint32
    endReached  bool
}

type Matcher interface {
    New(t *Lexer) Matcher
    Match(t *Lexer) (length uint32)
    Consume(t *Lexer, length uint32)
}

func NewScheme() *LexerScheme {
    return &LexerScheme{
        1,
        []Matcher{},
        END,
        []TokenTypeMetadata{
            {"a character that is not a valid token", "UNKNOWN"},
            {"to the end", "END"},
        },
    }
}

func (l *LexerScheme) RequireLookahead(n uint32) {
    if n > l.maxSafePeek {
        l.maxSafePeek = n
    }
}

func (l *LexerScheme) AddMatcher(matcher Matcher) {
    l.matchers = append(l.matchers, matcher)
}

func (l *LexerScheme) NewTokenType(metadata TokenTypeMetadata) TokenType {
    l.lastTokenType++
    l.metadata = append(l.metadata, metadata)

    return l.lastTokenType
}

func (l *LexerScheme) GetTokenTypeMetadata(tokenType TokenType) TokenTypeMetadata {
    return l.metadata[tokenType]
}

func (l *LexerScheme) Lex(source string) *Lexer {
    capacity := l.maxSafePeek

    job := &Lexer{
        make([]Matcher, len(l.matchers)),
        source,
        0,
        make([]Token, capacity),
        0,
        0,
        []Token{},
        l.maxSafePeek,
        false,
    }

    for i, matcher := range l.matchers {
        job.Matchers[i] = matcher.New(job)
    }

    job.fillTokenBuffer()

    return job
}

func (l *Lexer) Get(i uint32) (byte, bool) {
    offset := l.Position + i

    if offset >= uint32(len(l.Data)) {
        return 0, false
    }

    return l.Data[offset], true
}

func (l *Lexer) GetNextN(length uint32) string {
    return l.Data[l.Position:l.Position + length]
}

func (l *Lexer) Emit(token Token) {
    if l.write - l.read == uint32(len(l.buffer)) {
        l.spillage = append(l.spillage, token)

        return
    }

    l.buffer[l.write % uint32(len(l.buffer))] = token
    l.write++
}

func (l *Lexer) Peek(n uint32) Token {
    if n >= l.maxSafePeek {
        panic("attempt to peek more tokens in advance than expected")
    }

    read := l.read + n

    if read >= l.write {
        return Token{END, l.Data[len(l.Data):]}
    }

    return l.buffer[read % uint32(len(l.buffer))]
}

func (l *Lexer) Advance() {
    l.read++

    if l.read + l.maxSafePeek >= l.write {
        l.fillTokenBuffer()
    }
}

func (l *Lexer) fillTokenBuffer() {
    if len(l.spillage) > 0 {
        nEmpty := uint32(len(l.buffer)) - (l.write - l.read)
        nSpillage := uint32(len(l.spillage))

        nFill := min(nEmpty, nSpillage)

        for i := range nFill {
            l.buffer[l.write % uint32(len(l.buffer))] = l.spillage[i]
            l.write++
        }

        l.spillage = l.spillage[nFill:]
    }

    for !l.endReached && l.write - l.read != uint32(len(l.buffer)) {
        if l.Position >= uint32(len(l.Data)) {
            l.endReached = true

            return
        }

        largestLength := uint32(0)
        var matcherWithLargestLength Matcher = nil

        for _, matcher := range l.Matchers {
            length := matcher.Match(l)

            if length > largestLength {
                largestLength = length
                matcherWithLargestLength = matcher
            }
        }

        if matcherWithLargestLength != nil {
            matcherWithLargestLength.Consume(l, largestLength)
            l.Position += largestLength
        } else {
            l.Emit(Token{Type: UNKNOWN, Value: l.GetNextN(1)})

            l.Position++
        }
    }
}

func CheckTokens(t *testing.T, scheme *LexerScheme, expected []Token, text string) {
    actual := make([]Token, 0, len(expected))

    lexer := scheme.Lex(text)

    for current := lexer.Peek(0); current.Type != END; current = lexer.Peek(0) {
        actual = append(actual, current)
        lexer.Advance()
    }

    messenger := messenger.New()
    messenger.AddOutput(logrendering.New(os.Stdout))

    lexerDebugger := NewLexerDisplayer(scheme, os.Stdout, messenger)

    if len(expected) != len(actual) {
        lexerDebugger.DisplayTokensDiff(text, actual, expected)
        fmt.Println("")
        t.Fatal("Expected", len(expected), "tokens but got", len(actual), "tokens")
    }

    ok := true

    for i := range len(expected) {
        if actual[i].Type != expected[i].Type {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect type)",
            )

            ok = false

            break
        } else if actual[i].Value != expected[i].Value {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect value)",
            )

            ok = false

            break
        } else if unsafe.StringData(actual[i].Value) != unsafe.StringData(expected[i].Value) {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect string address)",
            )

            ok = false

            break
        }
    }

    if !ok {
        lexerDebugger.DisplayTokensDiff(text, actual, expected)
        fmt.Println()
    }
}
