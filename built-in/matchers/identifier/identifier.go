package identifiers

import (
	"minimal/minimal-lang/built-in/lexer"
)

type IdentifierMatcher struct {
	tokenType lexer.TokenType
}

func NewIdentifierMatcher(tt lexer.TokenType) *IdentifierMatcher {
	return &IdentifierMatcher{tt}
}

func (i *IdentifierMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return i
}

func (i *IdentifierMatcher) Match(l *lexer.Lexer) uint32 {
	firstChar, _ := l.Get(0)

	if !isAlphaOrUnicode(firstChar) {
		return 0
	}

	pos := uint32(1)

	for c, ok := l.Get(pos); ok && (isAlphaOrUnicode(c) || isDigit(c)); c, ok = l.Get(pos) {
		pos++
	}

	return pos
}

func (i *IdentifierMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: i.tokenType, Value: l.GetNextN(length)})
}

const asciiMax = 127

func isAlphaOrUnicode(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char > asciiMax
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
