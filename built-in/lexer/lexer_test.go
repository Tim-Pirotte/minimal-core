package lexer

import (
	"testing"
)

func TestLexEmpty(t *testing.T) {
	l := NewLexer()
	CheckTokens(t, l, []Token{}, "")
}

func TestLexUnknown(t *testing.T) {
	source := "a"

	l := NewLexer()
	expected := []Token{
		{Type: UNKNOWN, Value: source},
	}

	CheckTokens(t, l, expected, source)
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
		t.Emit(Token{s.tokenType, t.Data[:0]})
	}
}

func TestSpillage(t *testing.T) {
	source := "a"

	l := NewLexer()
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
	l := NewLexer()

	for b.Loop() {
		l.Lex(source, &testStopper{}, 1)
	}
}
