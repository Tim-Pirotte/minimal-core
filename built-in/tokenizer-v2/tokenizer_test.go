package tokenizerv2

import (
	"fmt"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/primitives"
	"os"
	"testing"
)

type testStopper struct {}

func (t *testStopper) End(s *TokenizerJob) bool {
	return s.Position >= uint(len(s.Data))
}

func CheckTokens(t *testing.T, tokenizer *Tokenizer, expected []Token, text string) {
	actual := make([]Token, 0, len(expected))

	tokenizerJob := tokenizer.Tokenize(text, &testStopper{}, 1)

	for current := tokenizerJob.Peek(0); current.Type != END; current = tokenizerJob.Peek(0) {
		actual = append(actual, current)
		tokenizerJob.Advance()
	}

	sourceGen := logging.GetTestLogSource(io.Discard)
	tokenizerDebugger := NewTokenizerDebugger(tokenizer, sourceGen, os.Stdout)

	if len(expected) != len(actual) {
		tokenizerDebugger.DisplayTokensDiff(actual, expected)
		fmt.Println("")
		t.Fatal("Expected", len(expected), "tokens but got", len(actual), "tokens")
	}

	ok := true

	for i := range len(expected) {
		if actual[i].Type != expected[i].Type {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]),
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect type)",
			)

			ok = false

			break
		} else if actual[i].Value != expected[i].Value {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]),
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect value)",
			)

			ok = false

			break
		} else if actual[i].Range != expected[i].Range {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]),
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect range)",
			)

			ok = false

			break
		}
	}

	if !ok {
		tokenizerDebugger.DisplayTokensDiff(actual, expected)
		fmt.Println("")
	}
}

func TestLexEmpty(t *testing.T) {
	tokenizer := NewTokenizer()
	CheckTokens(t, tokenizer, []Token{}, "")
}

func TestLexUnknown(t *testing.T) {
	tokenizer := NewTokenizer()
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	CheckTokens(t, tokenizer, expected, "a")
}

type spillageTestMatcher struct {
	tokenType TokenType
}

func newSpillageTestMatcher(tokenizer *Tokenizer) *spillageTestMatcher {
	return &spillageTestMatcher{
		tokenizer.NewTokenType(TokenTypeMetadata{"spillage", "Spillage"}),
	}
}

func (s *spillageTestMatcher) Match(t *TokenizerJob) uint {
	return 1
}

func (s *spillageTestMatcher) Consume(t *TokenizerJob, length uint) {
	for range 5 {
		t.Emit(Token{s.tokenType, "", primitives.Range{}})
	}
}

func TestSpillage(t *testing.T) {
	tokenizer := NewTokenizer()
	s := newSpillageTestMatcher(tokenizer)
	tokenizer.AddMatcher(s)

	expected := []Token{
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
	}

	CheckTokens(t, tokenizer, expected, "a")
}

func Benchmark(b *testing.B) {
	source := string(make([]byte, 1_000_000))
	tokenizer := NewTokenizer()

	for b.Loop() {
		tokenizer.Tokenize(source, &testStopper{}, 1)
	}
}
