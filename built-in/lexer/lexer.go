package lexer

import (
	"fmt"
	"minimal/minimal-core/built-in/messaging"
	"os"
	"testing"
	"unsafe"
)

type TokenType uint

const (
    UNKNOWN TokenType = iota
    // Always keep END at the end
    END
)

type Token struct {
    Type  TokenType
    Value string
}

type TokenTypeMetadata struct {
    DisplayName string
    DebugName   string
}

type Lexer struct {
    matchers      []Matcher
    lastTokenType TokenType
    metadata      []TokenTypeMetadata
}

type LexerJob struct {
    Matchers    []Matcher
    Data        string
    Position    uint
    buffer      []Token
    read        uint
    write       uint
    spillage    []Token
    minSafePeek uint
    endReached  bool
}

type Matcher interface {
    New(t *LexerJob) Matcher
    Match(t *LexerJob) (length uint)
    Consume(t *LexerJob, length uint)
}

func NewLexer() *Lexer {
    return &Lexer{
        []Matcher{},
        END,
        []TokenTypeMetadata{
            {"a character that is not a valid token", "UNKNOWN"},
            {"to the end", "END"},
        },
    }
}

func (l *Lexer) AddMatcher(matcher Matcher) {
    l.matchers = append(l.matchers, matcher)
}

func (l *Lexer) NewTokenType(metadata TokenTypeMetadata) TokenType {
    l.lastTokenType++
    l.metadata = append(l.metadata, metadata)

    return l.lastTokenType
}

func (l *Lexer) GetTokenTypeMetadata(tokenType TokenType) TokenTypeMetadata {
    return l.metadata[tokenType]
}

func (l *Lexer) Lex(source string, minSafePeek uint) *LexerJob {
    capacity := minSafePeek

    job := &LexerJob{
        make([]Matcher, len(l.matchers)),
        source,
        0,
        make([]Token, capacity),
        0,
        0,
        []Token{},
        minSafePeek,
        false,
    }

    for i, matcher := range l.matchers {
        job.Matchers[i] = matcher.New(job)
    }

    job.fillTokenBuffer()

    return job
}

func (l *LexerJob) Get(i uint) (byte, bool) {
    offset := l.Position + i

    if offset >= uint(len(l.Data)) {
        return 0, false
    }

    return l.Data[offset], true
}

func (l *LexerJob) GetNextN(length uint) string {
    return l.Data[l.Position:l.Position + length]
}

func (l *LexerJob) Emit(token Token) {
    if l.write - l.read == uint(len(l.buffer)) {
        l.spillage = append(l.spillage, token)

        return
    }

    l.buffer[l.write % uint(len(l.buffer))] = token
    l.write++
}

func (l *LexerJob) Peek(n uint) Token {
    if n >= l.minSafePeek {
        panic("attempt to peek more tokens in advance than expected")
    }

    read := l.read + n

    if read >= l.write {
        // TODO test this for off by one errors
        return Token{END, l.Data[len(l.Data):]}
    }

    return l.buffer[read % uint(len(l.buffer))]
}

func (l *LexerJob) Advance() {
    l.read++

    if l.read + l.minSafePeek >= l.write {
        l.fillTokenBuffer()
    }
}

func (l *LexerJob) fillTokenBuffer() {
    if len(l.spillage) > 0 {
        nEmpty := uint(len(l.buffer)) - (l.write - l.read)
        nSpillage := uint(len(l.spillage))

        nFill := min(nEmpty, nSpillage)

        for i := range nFill {
            l.buffer[l.write % uint(len(l.buffer))] = l.spillage[i]
            l.write++
        }

        l.spillage = l.spillage[nFill:]
    }

    // TODO check if we can put the end condition higher up
    for l.write - l.read != uint(len(l.buffer)) && !l.endReached {
        if l.Position >= uint(len(l.Data)) {
            l.endReached = true

            return
        }

        largestLength := uint(0)
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

func CheckTokens(t *testing.T, lexer *Lexer, expected []Token, text string) {
    actual := make([]Token, 0, len(expected))

    lexerJob := lexer.Lex(text, 1)

    for current := lexerJob.Peek(0); current.Type != END; current = lexerJob.Peek(0) {
        actual = append(actual, current)
        lexerJob.Advance()
    }

    // TODO What with this messenger?
    lexerDebugger := NewLexerDebugger(lexer, os.Stdout, messaging.NewMessenger())

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
