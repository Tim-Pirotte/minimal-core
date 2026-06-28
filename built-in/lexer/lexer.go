package lexer

import (
	"fmt"
	"io"
	"math"
	logging "minimal/minimal-core/built-in/internal-logging"
	"os"
	"testing"
	"unsafe"
)

type TokenType uint

const (
	UNKNOWN TokenType = math.MaxUint - iota
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
	tokenTypesMetadata map[TokenType]TokenTypeMetadata
}

type LexerJob struct {
	matchers    []Matcher
	Data        string
	Position    uint
	buffer      []Token
	read        uint
	write       uint
	spillage    []Token
	minSafePeek uint
	stopper     Stopper
	endReached  bool
}

type Matcher interface {
	New(t *LexerJob) Matcher
	Match(t *LexerJob) (length uint)
	Consume(t *LexerJob, length uint)
}

type Stopper interface {
	End(t *LexerJob) bool
}

func NewLexer() *Lexer {
	return &Lexer{
		[]Matcher{},
		TokenType(0),
		map[TokenType]TokenTypeMetadata{
			UNKNOWN: {"a character that is not a valid token", "UNKNOWN"},
			END: {"to the end", "END"},
		},
	}
}

func (l *Lexer) AddMatcher(matcher Matcher) {
	l.matchers = append(l.matchers, matcher)
}

func (l *Lexer) NewTokenType(metadata TokenTypeMetadata) TokenType {
	l.lastTokenType++
	l.tokenTypesMetadata[l.lastTokenType] = metadata

	return l.lastTokenType
}

func (l *Lexer) GetTokenTypeMetadata(tokenType TokenType) (TokenTypeMetadata, bool) {
	v, ok := l.tokenTypesMetadata[tokenType]

	return v, ok
}

func (l *Lexer) Lex(source string, stopper Stopper, minSafePeek uint) *LexerJob {
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
		stopper,
		false,
	}

	for i, matcher := range l.matchers {
		job.matchers[i] = matcher.New(job)
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

	// TODO separate bounds check from end function and
	// check if we can put the end condition higher up
	for l.write - l.read != uint(len(l.buffer)) && !l.endReached {
		if l.stopper.End(l) {
			l.endReached = true

			return
		}

		largestLength := uint(0)
		var matcherWithLargestLength Matcher = nil

		for _, matcher := range l.matchers {
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

type testStopper struct {}

func (t *testStopper) End(s *LexerJob) bool {
	return s.Position >= uint(len(s.Data))
}

func CheckTokens(t *testing.T, lexer *Lexer, expected []Token, text string) {
	actual := make([]Token, 0, len(expected))

	lexerJob := lexer.Lex(text, &testStopper{}, 1)

	for current := lexerJob.Peek(0); current.Type != END; current = lexerJob.Peek(0) {
		actual = append(actual, current)
		lexerJob.Advance()
	}

	sourceGen := logging.GetTestLogSource(io.Discard)
	lexerDebugger := NewLexerDebugger(lexer, sourceGen, os.Stdout)

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
