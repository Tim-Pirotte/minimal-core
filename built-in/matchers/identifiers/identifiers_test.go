package identifiers

import (
	"fmt"
	"io"
	tokenizerdebugger "minimal/minimal-core/built-in/debug/tokenizer-debugger"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/primitives"
	eofstopper "minimal/minimal-core/built-in/stoppers/eof-stopper"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"os"
	"strings"
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

func TestLexIdentifier(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
	}

	actual := tokenizer.Tokenize("identifier1", eofstopper.NewEOFStopper())

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexMultipleIdentifiers(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "identifier1", Range: primitives.Range{Start: 0, Length: 11}},
		{Type: tokenizerv2.UNKNOWN, Value: " ", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Range: primitives.Range{Start: 12, Length: 11}},
	}

	actual := tokenizer.Tokenize("identifier1 identifier2", eofstopper.NewEOFStopper())

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexZeroWidthJoiner(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
	}

	actual := tokenizer.Tokenize("🐻‍❄️", eofstopper.NewEOFStopper())

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexUnicode(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: identifierType, Value: "🐻‍❄️", Range: primitives.Range{Start: 0, Length: 13}},
	}

	actual := tokenizer.Tokenize("🐻‍❄️", eofstopper.NewEOFStopper())

	testTokensEqual(t, &tokenizer, expected, actual)
}

func TestLexStartingWithNumber(t *testing.T) {
	tokenizer, identifierType := getTokenizer()

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Range: primitives.Range{Start: 1, Length: 10}},
	}

	actual := tokenizer.Tokenize("1identifier", eofstopper.NewEOFStopper())

	testTokensEqual(t, &tokenizer, expected, actual)
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
		
		tokens := tokenizer.Tokenize(input, eofstopper.NewEOFStopper())

		if len(tokens) != 1 {
			t.Fatalf("expected 1 token, got %d tokens", len(tokens))
		}

		if tokens[0].Type != identifierType {
			t.Fatalf("expected identifier token, got %v", tokens[0].Type)
		}

		if tokens[0].Value != input {
			t.Fatalf("expected %q, got %q", input, tokens[0].Value)
		}

		if tokens[0].Range.Start != 0 {
			t.Fatalf("expected start=0, got %d", tokens[0].Range.Start)
		}

		if tokens[0].Range.Length != uint(len(input)) {
			t.Fatalf(
				"expected length=%d, got %d",
				len(input),
				tokens[0].Range.Length,
			)
		}
	})
}
