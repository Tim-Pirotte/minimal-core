package whitespace

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"testing"
)

func TestLexSpace(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierMatcher := NewWhiteSpaceMatcher()
	tokenizer.AddMatcher(identifierMatcher)

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "a", Range: primitives.Range{Start: 12, Length: 1}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, " \t \t        a\t            ")
}
