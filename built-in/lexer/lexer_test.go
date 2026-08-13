package lexer

import (
	"testing"
)

// TODO tests for the printer

func TestLexEmpty(t *testing.T) {
	l := New(1)
	CheckTokens(t, l, []Token{}, "")
}

func TestLexUnknown(t *testing.T) {
	source := "a"

	l := New(1)
	expected := []Token{
		{Type: UNKNOWN, Value: source},
	}

	CheckTokens(t, l, expected, source)
}

type spillageTestMatcher struct {
	tokenType TokenType
}

func newSpillageTestMatcher(lexer *LexerScheme) *spillageTestMatcher {
	return &spillageTestMatcher{
		lexer.NewTokenType(TokenTypeMetadata{"spillage", "Spillage"}),
	}
}

func (s *spillageTestMatcher) New(_ *Lexer) Matcher {
	return s
}

func (s *spillageTestMatcher) Match(t *Lexer) uint {
	return 1
}

func (s *spillageTestMatcher) Consume(t *Lexer, length uint) {
	for range 5 {
		t.Emit(Token{s.tokenType, t.Data[:0]})
	}
}

func TestSpillage(t *testing.T) {
	source := "a"

	l := New(1)
	s := newSpillageTestMatcher(l)
	l.AddMatcher(s)

	expected := []Token{
		{s.tokenType, source[:0]},
		{s.tokenType, source[:0]},
		{s.tokenType, source[:0]},
		{s.tokenType, source[:0]},
		{s.tokenType, source[:0]},
	}

	CheckTokens(t, l, expected, source)
}

func Benchmark(b *testing.B) {
	source := string(make([]byte, 1_000_000))
	l := New(1)

	for b.Loop() {
		l.Lex(source)
	}
}
