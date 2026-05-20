package identifiers

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"testing"
)

func testTokensEqual(t *testing.T, expected, actual []tokenizerv2.Token) {
	if len(expected) != len(actual) {
		t.Fatal("Expected", len(expected), "tokens but got", len(actual), "tokens", "\nExpected:", expected, "\nActual:", actual)
	}

	for i := range len(expected) {
		if actual[i].Type != expected[i].Type {
			t.Error("Expected", expected[i], "but got", actual[i], "(incorrect type)")
		} else if actual[i].Value != expected[i].Value {
			t.Error("Expected", expected[i], "but got", actual[i], "(incorrect value)")
		} else if actual[i].Range != expected[i].Range {
			t.Error("Expected", expected[i], "but got", actual[i], "(incorrect range)")
		}
	}
}

func TestLexIdentifiers(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	tokenizer.AddMatcher(&identifierMatcher)

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 11, Length: 0}},
	}

	actual := tokenizer.Tokenize("identifier1")

	testTokensEqual(t, expected, actual)
}

func TestLexMultipleIdentifiers(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	tokenizer.AddMatcher(&identifierMatcher)

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.UNKNOWN, Value: " ", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Range: primitives.Range{Start: 12, Length: 11}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 23, Length: 0}},
	}

	actual := tokenizer.Tokenize("identifier1 identifier2")

	testTokensEqual(t, expected, actual)
}

func TestLexUnicode(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	tokenizer.AddMatcher(&identifierMatcher)

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 13, Length: 0}},
	}

	actual := tokenizer.Tokenize("🐻‍❄️")

	testTokensEqual(t, expected, actual)
}

func TestLexStartingWithNumber(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	identifierType := tokenizer.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	tokenizer.AddMatcher(&identifierMatcher)

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Range: primitives.Range{Start: 1, Length: 10}},
		{Type: tokenizerv2.EOF, Value: "", Range: primitives.Range{Start: 11, Length: 0}},
	}

	actual := tokenizer.Tokenize("1identifier")

	testTokensEqual(t, expected, actual)
}
