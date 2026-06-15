package lexer

import (
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	tokenizer := NewLexer()
	CheckTokens(t, tokenizer, []Token{}, "")
}

func TestLexUnknown(t *testing.T) {
	tokenizer := NewLexer()
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	CheckTokens(t, tokenizer, expected, "a")
}

type spillageTestMatcher struct {
	tokenType TokenType
}

func newSpillageTestMatcher(tokenizer *Lexer) *spillageTestMatcher {
	return &spillageTestMatcher{
		tokenizer.NewTokenType(TokenTypeMetadata{"spillage", "Spillage"}),
	}
}

func (s *spillageTestMatcher) Match(t *LexerJob) uint {
	return 1
}

func (s *spillageTestMatcher) Consume(t *LexerJob, length uint) {
	for range 5 {
		t.Emit(Token{s.tokenType, "", primitives.Range{}})
	}
}

func TestSpillage(t *testing.T) {
	tokenizer := NewLexer()
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
	tokenizer := NewLexer()

	for b.Loop() {
		tokenizer.Lex(source, &testStopper{}, 1)
	}
}
