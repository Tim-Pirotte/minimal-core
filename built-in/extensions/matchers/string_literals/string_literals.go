package stringliterals

import (
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
	"strings"
)

type StringLiteralMatcher struct {
	tokenType domain.TokenType
}

func NewStringLiteralMatcher(tt domain.TokenType) StringLiteralMatcher {
	return StringLiteralMatcher{tt}
}

func (s *StringLiteralMatcher) Match(so *tokenizer.Source) (uint, domain.TokenType, string) {
	firstChar, _ := so.Get(0)

	if !(firstChar == '"') {
		return 0, domain.IGNORE, ""
	}

	sb := strings.Builder{}

	pos := 1
	for {
		currentChar, ok := so.Get(pos)

		if !ok {
			// TODO log error
			panic("No closing character \"")
		}

		if currentChar == '"' {
			pos++
			break;
		}

		if currentChar == '\\' {
			sb.WriteByte(currentChar)
			// The characters that can close a string are all one byte long
			charToSkip, ok := so.Get(pos + 1)

			if !ok {
				// TODO log error
				panic("Nothing to escape")
			}

			sb.WriteByte(charToSkip)

			pos += 2
		} else {
			sb.WriteByte(currentChar)
			pos++
		}
	}

	return uint(pos), s.tokenType, sb.String()
}
