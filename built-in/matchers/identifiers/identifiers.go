package identifiers

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type IdentifierMatcher struct {
	tokenType tokenizerv2.TokenType
}

func NewIdentifierMatcher(tt tokenizerv2.TokenType) *IdentifierMatcher {
	return &IdentifierMatcher{tt}
}

func (i *IdentifierMatcher) Match(t *tokenizerv2.TokenizerJob) uint {
	firstChar, ok := t.Get(0)

	if !ok || !isAlphaOrUnicode(firstChar) {
		return 0
	}

	pos := uint(1)

	for char, ok := t.Get(pos); ok && isValidIdentifierChar(char); char, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (i *IdentifierMatcher) Consume(t *tokenizerv2.TokenizerJob, length uint) {
	identifier, _ := t.GetRange(t.Position, length)

	t.Emit(tokenizerv2.Token{
		Type: i.tokenType,
		Value: identifier,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}

func isValidIdentifierChar(char byte) bool {
	return isAlphaOrUnicode(char) || isDigit(char)
}

const asciiMax = 127

func isAlphaOrUnicode(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char > asciiMax
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
