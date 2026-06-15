package numbers

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType) {
	tokenizer := lexer.NewLexer()
	numberType := tokenizer.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a number", DebugName: "Number"},
	)

	identifierMatcher := NewNumberMatcher(numberType)
	tokenizer.AddMatcher(identifierMatcher)

	return tokenizer, numberType
}

func TestNumber(t *testing.T) {
	tokenizer, numberType := getLexer()

	expected := []lexer.Token{
		{Type: numberType, Value: "0123456789", Range: primitives.Range{Start: 0, Length: 10}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "0123456789")
}

func TestMixed(t *testing.T) {
	tokenizer, numberType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "/", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: numberType, Value: "0", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: ":", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: numberType, Value: "12345", Range: primitives.Range{Start: 3, Length: 5}},
		{Type: lexer.UNKNOWN, Value: " ", Range: primitives.Range{Start: 8, Length: 1}},
		{Type: numberType, Value: "6789", Range: primitives.Range{Start: 9, Length: 4}},
	}

	// / and : are next to 0 and 9 respectively in the ASCII table
	lexer.CheckTokens(t, tokenizer, expected, "/0:12345 6789")
}
