package identifiers

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

type IdentifierMatcher struct {
	tokenType lexer.TokenType
}

func NewIdentifierMatcher(tt lexer.TokenType) *IdentifierMatcher {
	return &IdentifierMatcher{tt}
}

func (i *IdentifierMatcher) Match(t *lexer.LexerJob) uint {
	firstChar, ok := t.Get(0)

	if !ok || !isAlphaOrUnicode(firstChar) {
		return 0
	}

	pos := uint(1)

	for c, ok := t.Get(pos); ok && (isAlphaOrUnicode(c) || isDigit(c)); c, ok = t.Get(pos) {
		pos++
	}

	return pos
}

func (i *IdentifierMatcher) Consume(t *lexer.LexerJob, length uint) {
	identifier, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type: i.tokenType,
		Value: identifier,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}

const asciiMax = 127

func isAlphaOrUnicode(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char > asciiMax
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
