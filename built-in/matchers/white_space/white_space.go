package whitespace

import (
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type WhiteSpaceMatcher struct{}

func NewWhiteSpaceMatcher() *WhiteSpaceMatcher {
	return &WhiteSpaceMatcher{}
}

func (*WhiteSpaceMatcher) Match(t *tokenizerv2.TokenizerJob) uint {
	pos := uint(0)

	for ch, ok := t.Get(pos); ok && isWhiteSpace(ch); ch, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (*WhiteSpaceMatcher) Consume(t *tokenizerv2.TokenizerJob, length uint) {}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t'
}
