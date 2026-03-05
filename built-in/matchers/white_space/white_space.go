package whitespace

import (
	"minimal/minimal-core/built-in/tokenizer"
)

type WhiteSpaceMatcher struct{}

func NewWhiteSpaceMatcher() WhiteSpaceMatcher {
	return WhiteSpaceMatcher{}
}

func (*WhiteSpaceMatcher) Match(s *tokenizer.Source) (uint, tokenizer.TokenType, string) {
	pos := 0
	
	for {
		ch, ok := s.Get(pos)
		
		if !ok || !isWhiteSpace(ch) {
			break
		}

		pos++
	}

	return uint(pos), tokenizer.IGNORE, ""
}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t'
}
