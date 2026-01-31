package whitespace

import (
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
)

type WhiteSpaceMatcher struct{}

func NewWhiteSpaceMatcher() WhiteSpaceMatcher {
	return WhiteSpaceMatcher{}
}

func (*WhiteSpaceMatcher) Match(s *tokenizer.Source) (bool, uint, domain.TokenType, string) {
	pos := 0
	
	for {
		ch, ok := s.Get(pos)
		
		if !ok || !isWhiteSpace(ch) {
			break
		}

		pos++
	}

	return true, uint(pos), domain.IGNORE, ""
}

func isWhiteSpace(b byte) bool {
	return b == ' ' || b == '\t'
}
