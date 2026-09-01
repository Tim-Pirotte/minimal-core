package whitespace

import (
	"minimal/minimal-core/built-in/lexer"
)

type SpaceMatcher struct{}

func NewSpaceMatcher() *SpaceMatcher {
	return &SpaceMatcher{}
}

func (s *SpaceMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return s
}

func (*SpaceMatcher) Match(l *lexer.Lexer) uint32 {
	pos := uint32(0)

	for c, ok := l.Get(pos); ok && c == ' '; c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (*SpaceMatcher) Consume(_ *lexer.Lexer, _ uint32) {}
