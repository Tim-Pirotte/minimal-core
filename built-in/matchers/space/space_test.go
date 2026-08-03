package whitespace

import (
	"minimal/minimal-core/built-in/lexer"
	"testing"
)

func TestLexSpace(t *testing.T) {
	source := " \t        a            "

	l := lexer.New(1)
	identifierMatcher := NewSpaceMatcher()
	l.AddMatcher(identifierMatcher)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[1:2]},
		{Type: lexer.UNKNOWN, Value: source[10:11]},
	}

	lexer.CheckTokens(t, l, expected, source)
}
