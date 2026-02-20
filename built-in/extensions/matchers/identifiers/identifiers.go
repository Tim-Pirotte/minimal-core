package identifiers

import (
	"minimal/minimal-core/built-in/tokenizer"
	"strings"
)

type IdentifierMatcher struct {
	tokenType tokenizer.TokenType
}

func NewIdentifierMatcher(tt tokenizer.TokenType) IdentifierMatcher {
	return IdentifierMatcher{tt}
}

func (i *IdentifierMatcher) Match(s *tokenizer.Source) (uint, tokenizer.TokenType, string) {
	firstChar, _ := s.Get(0)

	if !isAlphaOrUnicode(firstChar) {
		return 0, tokenizer.IGNORE, ""
	}

	sb := strings.Builder{}
	sb.WriteByte(firstChar)

	pos := 1
	for {
		currentChar, ok := s.Get(pos)

		if !ok || (!isAlphaOrUnicode(currentChar) && !isDigit(currentChar)) {
			break;
		}

		sb.WriteByte(currentChar)

		pos++
	}

	return uint(pos), i.tokenType, sb.String()
}

const asciiMax = 127

func isAlphaOrUnicode(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char > asciiMax
}

func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}
