package identifiers

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type IdentifierMatcher struct {
	tokenType tokenizerv2.TokenType
}

func NewIdentifierMatcher(tt tokenizerv2.TokenType) IdentifierMatcher {
	return IdentifierMatcher{tt}
}

func (i *IdentifierMatcher) Match(s *tokenizerv2.TokenizerState) uint {
	firstChar, ok := s.Get(0)

	if !ok || !isAlphaOrUnicode(firstChar) {
		return 0
	}

	pos := uint(1)

	for char, ok := s.Get(pos); ok && isValidIdentifierChar(char); char, ok = s.Get(pos) {
		pos++
	}

	return pos
}

func (i *IdentifierMatcher) Consume(s *tokenizerv2.TokenizerState, length uint) {
	identifier, _ := s.GetRange(s.Position, length)

	s.Emit(tokenizerv2.Token{
		Type: i.tokenType, 
		Value: identifier, 
		Range: primitives.Range{Start: s.Position, Length: length}},
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
