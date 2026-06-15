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

	for ch, ok := t.Get(pos); ok && isWhiteSpace(ch); ch, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (*SpaceMatcher) Consume(t *lexer.LexerJob, length uint) {}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t'
}
