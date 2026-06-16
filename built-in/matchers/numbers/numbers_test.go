package numbers

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType) {
	l := lexer.NewLexer()
	numberType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a number", DebugName: "Number"},
	)

	identifierMatcher := NewNumberMatcher(numberType)
	l.AddMatcher(identifierMatcher)

	return l, numberType
}

func TestNumber(t *testing.T) {
	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: "0123456789", Range: primitives.Range{Start: 0, Length: 10}},
	}

	lexer.CheckTokens(t, l, expected, "0123456789")
}

func TestMixed(t *testing.T) {
	l, numberType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "/", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: numberType, Value: "0", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: ":", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: numberType, Value: "12345", Range: primitives.Range{Start: 3, Length: 5}},
		{Type: lexer.UNKNOWN, Value: " ", Range: primitives.Range{Start: 8, Length: 1}},
		{Type: numberType, Value: "6789", Range: primitives.Range{Start: 9, Length: 4}},
	}

	// / and : are next to 0 and 9 respectively in the ASCII table
	lexer.CheckTokens(t, l, expected, "/0:12345 6789")
}
