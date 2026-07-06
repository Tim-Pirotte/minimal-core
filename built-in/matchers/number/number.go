package numbers

import (
	"minimal/minimal-core/built-in/lexer"
)

type NumberMatcher struct {
	tokenType lexer.TokenType
}

func NewNumberMatcher(tt lexer.TokenType) *NumberMatcher {
	return &NumberMatcher{tt}
}

func (n *NumberMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
	return n
}

func (n *NumberMatcher) Match(l *lexer.LexerJob) uint {
	pos := uint(0)

	for c, ok := l.Get(pos); ok && '0' <= c && c <= '9'; c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *NumberMatcher) Consume(l *lexer.LexerJob, length uint) {
	l.Emit(lexer.Token{Type: i.tokenType, Value: l.GetNextN(length)})
}
