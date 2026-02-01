package stringliterals

import (
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
	"strings"
)

type StringLiteralMatcher struct {
	stringTt             domain.TokenType
	interpolatedOpenTt   domain.TokenType
	interpolatedMiddleTt domain.TokenType
	interpolatedCloseTt  domain.TokenType

	braceCounter         uint
}

func NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt domain.TokenType) StringLiteralMatcher {
	return StringLiteralMatcher{stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt, 0}
}

func (s *StringLiteralMatcher) Match(so *tokenizer.Source) (uint, domain.TokenType, string) {
	firstChar, _ := so.Get(0)

	var startsWithQuotation bool

	switch firstChar {
	case '{':
		if s.braceCounter > 0 {
			s.braceCounter++
		}

		return 0, domain.IGNORE, ""
	case '}':
		if s.braceCounter > 0 {
			s.braceCounter--
			startsWithQuotation = false

			if s.braceCounter != 0 {
				return 0, domain.IGNORE, ""
			}
		} else {
			return 0, domain.IGNORE, ""
		}
	case '"':
		startsWithQuotation = true
	default:
		return 0, domain.IGNORE, ""
	}

	sb := strings.Builder{}

	pos := 1
	for {
		currentChar, ok := so.Get(pos)

		if !ok {
			if s.braceCounter > 0 {
				// TODO log error
				panic("No closing character }")
			} else {
				// TODO log error
				panic("No closing character \"")
			}
		}

		switch currentChar {
		case '"':
			if startsWithQuotation {
				return uint(pos) + 1, s.stringTt, sb.String()
			} else {
				return uint(pos) + 1, s.interpolatedCloseTt, sb.String()
			}
		case '{':
			s.braceCounter++

			if startsWithQuotation {
				return uint(pos) + 1, s.interpolatedOpenTt, sb.String()
			} else {
				return uint(pos) + 1, s.interpolatedMiddleTt, sb.String()
			}
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
}
