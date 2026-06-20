package lexer

import (
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	l := NewLexer()
	CheckTokens(t, l, []Token{}, "")
}

func TestLexUnknown(t *testing.T) {
	l := NewLexer()
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	CheckTokens(t, l, expected, "a")
}

type spillageTestMatcher struct {
	tokenType TokenType
}

func newSpillageTestMatcher(lexer *Lexer) *spillageTestMatcher {
	return &spillageTestMatcher{
		lexer.NewTokenType(TokenTypeMetadata{"spillage", "Spillage"}),
	}
}

func (s *spillageTestMatcher) New(_ *LexerJob) Matcher {
	return s
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
	l := NewLexer()
	s := newSpillageTestMatcher(l)
	l.AddMatcher(s)

	expected := []Token{
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
		{s.tokenType, "", primitives.Range{}},
	}

	CheckTokens(t, l, expected, "a")
}

func Benchmark(b *testing.B) {
	source := string(make([]byte, 1_000_000))
	l := NewLexer()

	for b.Loop() {
		l.Lex(source, &testStopper{}, 1)
	}
}
