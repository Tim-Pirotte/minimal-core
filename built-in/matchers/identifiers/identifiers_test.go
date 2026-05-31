package identifiers

import (
	"fmt"
	"io"
	tokenizerdebugger "minimal/minimal-core/built-in/debug/tokenizer-debugger"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"os"
	"testing"
)

func testTokensEqual(t *testing.T, tokenizer *tokenizerv2.Tokenizer, expected, actual []tokenizerv2.Token) {
	sourceGen := logging.GetTestLogSource(io.Discard)
	tokenizerDebugger := tokenizerdebugger.NewTokenizerDebugger(tokenizer, sourceGen, os.Stdout)

	if len(expected) != len(actual) {
		tokenizerDebugger.DisplayTokensDiff(actual, expected)
		fmt.Println("")
		t.Fatal("Expected", len(expected), "tokens but got", len(actual), "tokens")
	}

	ok := true

	for i := range len(expected) {
		if actual[i].Type != expected[i].Type {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]), 
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect type)",
			)

			ok = false

			break
		} else if actual[i].Value != expected[i].Value {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]), 
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect value)",
			)

			ok = false

			break
		} else if actual[i].Range != expected[i].Range {
			t.Error(
				"\nExpected\n", tokenizerDebugger.StringifyToken(expected[i]), 
				"\nbut got\n", tokenizerDebugger.StringifyToken(actual[i]), "(incorrect range)",
			)

			ok = false

			break
		}
	}

	if !ok {
		tokenizerDebugger.DisplayTokensDiff(actual, expected)
		fmt.Println("")
	}
}

func getTokenizer() (tokenizerv2.Tokenizer, tokenizerv2.TokenType) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{DisplayName: "an identifier", DebugName: "Identifier"})
	identifierMatcher := NewIdentifierMatcher(identifierType)

	tokenizer.AddMatcher(&identifierMatcher)

	return tokenizer, identifierType
}

func TestLexIdentifiers(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 11, Length: 0}},
	}

	actual := tokenizer.Tokenize("identifier1")

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexMultipleIdentifiers(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.UNKNOWN, Value: " ", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Range: primitives.Range{Start: 12, Length: 11}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 23, Length: 0}},
	}

	actual := tokenizer.Tokenize("identifier1 identifier2")

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexZeroWidthJoiner(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 13, Length: 0}},
	}

	actual := tokenizer.Tokenize("🐻‍❄️")

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexUnicode(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 13, Length: 0}},
	}

	actual := tokenizer.Tokenize("🐻‍❄️")

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexStartingWithNumber(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Range: primitives.Range{Start: 1, Length: 10}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 11, Length: 0}},
	}

	actual := tokenizer.Tokenize("1identifier")

	testTokensEqual(t, &tokenizer, expected, actual)
}
