package eol

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type EOLMatcher struct{
	tokenType lexer.TokenType
}

func NewEOLMatcher(tokenType lexer.TokenType) *EOLMatcher {
	return &EOLMatcher{}
}

func (*EOLMatcher) Match(t *lexer.LexerJob) uint {
	pos := uint(0)

	for c, ok := t.Get(pos); ok && (c == '\n' || c == '\r'); c, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (e *EOLMatcher) Consume(t *lexer.LexerJob, length uint) {
	endOfLine, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type: e.tokenType,
		Value: endOfLine,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
