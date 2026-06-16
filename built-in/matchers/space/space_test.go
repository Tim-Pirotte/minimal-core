package whitespace

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func TestLexSpace(t *testing.T) {
	l := lexer.NewLexer()
	identifierMatcher := NewSpaceMatcher()
	l.AddMatcher(identifierMatcher)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "\t", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 10, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, " \t        a            ")
}
