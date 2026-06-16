package identifiers

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"strings"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType) {
	tokenizer := lexer.NewLexer()
	identifierType := tokenizer.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "an identifier", DebugName: "Identifier"},
	)

	identifierMatcher := NewIdentifierMatcher(identifierType)
	tokenizer.AddMatcher(identifierMatcher)

	return tokenizer, identifierType
}

func TestIdentifier(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "identifier1")
}

func TestMultipleIdentifiers(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: lexer.UNKNOWN, Value: " ", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Range: primitives.Range{Start: 12, Length: 11}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "identifier1 identifier2")
}

func TestUnicode(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: "👾", Range: primitives.Range{Start: 0, Length: 4}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "👾")
}

func TestZeroWidthJoiner(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "🐻‍❄️")
}

func TestStartingWithNumber(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Range: primitives.Range{Start: 1, Length: 10}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "1identifier")
}

func TestBounds(t *testing.T) {
	tokenizer, identifierType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "`", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "az", Range: primitives.Range{Start: 1, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "{", Range: primitives.Range{Start: 3, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "@", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: identifierType, Value: "AZ", Range: primitives.Range{Start: 5, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "[", Range: primitives.Range{Start: 7, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "\u007f", Range: primitives.Range{Start: 8, Length: 1}},
		{Type: identifierType, Value: "\u0080\U0010ffff", Range: primitives.Range{Start: 9, Length: 6}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "`az{@AZ[\u007f\u0080\U0010ffff")
}

func FuzzLexIdentifier(f *testing.F) {
	tokenizer, identifierType := getLexer()

	f.Add("identifier")

	f.Fuzz(func(t *testing.T, input string) {
		if len(input) > 0 && '0' <= input[0] && input[0] <= '9' {
			input = input[1:]
		}

		var cleanedInput strings.Builder

		for _, c := range []byte(input) {
			if isAlphaOrUnicode(c) {
				cleanedInput.WriteString(string([]byte{c}))
			}
		}

		input = cleanedInput.String()

		if input == "" {
			return
		}

		expected := []lexer.Token{
			{
				Type: identifierType,
				Value: input,
				Range: primitives.Range{Start: 0, Length: uint(len(input))},
			},
		}

		lexer.CheckTokens(t, tokenizer, expected, input)
	})
}
