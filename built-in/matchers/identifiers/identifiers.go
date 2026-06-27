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

func (i *IdentifierMatcher) New(_ *lexer.LexerJob) lexer.Matcher {
	return i
}

func (i *IdentifierMatcher) Match(l *lexer.LexerJob) uint {
	firstChar, ok := l.Get(0)

	if !ok || !isAlphaOrUnicode(firstChar) {
		return 0
	}

	pos := uint(1)

	for c, ok := l.Get(pos); ok && (isAlphaOrUnicode(c) || isDigit(c)); c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *IdentifierMatcher) Consume(l *lexer.LexerJob, length uint) {
	l.Emit(lexer.Token{
		Type: i.tokenType,
		Value: l.GetNextN(length),
		Range: primitives.Range{Start: l.Position, Length: length}},
	)
}

const asciiMax = 127

func isAlphaOrUnicode(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char > asciiMax
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
