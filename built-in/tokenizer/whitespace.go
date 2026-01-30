package tokenizer

import "minimal/minimal-core/domain"

type WhiteSpaceMatcher struct{}

func NewWhiteSpaceMatcher() WhiteSpaceMatcher {
	return WhiteSpaceMatcher{}
}

func (*WhiteSpaceMatcher) Match(t *TokenizerConfig) (bool, uint, domain.TokenType, string) {
	pos := 0
	
	for {
		ch, ok := t.Get(pos)
		
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
