package numbers

import (
	"minimal/minimal-lang/built-in/lexer"
)

type NumberMatcher struct {
	tokenType lexer.TokenType
}

func NewNumberMatcher(tt lexer.TokenType) *NumberMatcher {
	return &NumberMatcher{tt}
}

func (n *NumberMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return n
}

func (n *NumberMatcher) Match(l *lexer.Lexer) uint32 {
	pos := uint32(0)

	for c, ok := l.Get(pos); ok && '0' <= c && c <= '9'; c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *NumberMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: i.tokenType, Value: l.GetNextN(length)})
}
