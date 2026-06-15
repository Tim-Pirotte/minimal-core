package numbers

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type NumberMatcher struct {
	tokenType lexer.TokenType
}

func NewNumberMatcher(tt lexer.TokenType) *NumberMatcher {
	return &NumberMatcher{tt}
}

func (n *NumberMatcher) Match(t *lexer.LexerJob) uint {
	pos := uint(0)

	for char, ok := t.Get(pos); ok && '0' <= char && char <= '9'; char, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (i *NumberMatcher) Consume(t *lexer.LexerJob, length uint) {
	number, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type: i.tokenType,
		Value: number,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
