package whitespace

import (
	"minimal/minimal-core/built-in/lexer"
)

type SpaceMatcher struct{}

func NewSpaceMatcher() *SpaceMatcher {
	return &SpaceMatcher{}
}

func (*SpaceMatcher) Match(t *lexer.LexerJob) uint {
	pos := uint(0)

	for c, ok := t.Get(pos); ok && c == ' '; c, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (*SpaceMatcher) Consume(_ *lexer.LexerJob, _ uint) {}
