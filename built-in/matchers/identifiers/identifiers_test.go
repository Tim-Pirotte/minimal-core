package identifiers

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"strings"
	"testing"
)

func getTokenizer() (*tokenizerv2.Tokenizer, tokenizerv2.TokenType) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType(
		tokenizerv2.TokenTypeMetadata{DisplayName: "an identifier", DebugName: "Identifier"},
	)

	identifierMatcher := NewIdentifierMatcher(identifierType)
	tokenizer.AddMatcher(identifierMatcher)

	return tokenizer, identifierType
}

func TestLexIdentifier(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "identifier1")
}

func TestLexMultipleIdentifiers(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.UNKNOWN, Value: " ", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Range: primitives.Range{Start: 12, Length: 11}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "identifier1 identifier2")
}

func TestLexUnicode(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "👾", Range: primitives.Range{Start: 0, Length: 4}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "👾")
}

func TestLexZeroWidthJoiner(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "🐻‍❄️")
}

func TestLexStartingWithNumber(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Range: primitives.Range{Start: 1, Length: 10}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "1identifier")
}

func FuzzLexIdentifier(f *testing.F) {
	tokenizer, identifierType := getTokenizer()

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

		expected := []tokenizerv2.Token{
			{
				Type: identifierType,
				Value: input,
				Range: primitives.Range{Start: 0, Length: uint(len(input))},
			},
		}

		tokenizerv2.CheckTokens(t, tokenizer, expected, input)
	})
}
