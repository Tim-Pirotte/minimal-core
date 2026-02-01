package stringliterals

import (
	"fmt"
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
	"reflect"
	"testing"
)

func TestLexStringLiteral(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringLiteralType := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringLiteralType, 9999, 9999, 9999)

	config.AddMatcher(&stringLiteralMatcher)

	expected := []domain.Token{
		{Type: stringLiteralType, Value: "Hello, world!", Span: domain.Span{Start: 0, Length: 15}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 15, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("\"Hello, world!\""))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i+1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
	}
}

func TestLexUnclosed(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringLiteralType := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringLiteralType, 9999, 9999, 9999)

	config.AddMatcher(&stringLiteralMatcher)

	// TODO check if it logs properly
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenizer.NewTokenizer(config, []byte("\"Hello, world!"))
}

func TestLexInterpolatedString(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	expected := []domain.Token{
		{Type: interpolatedOpenTt, Value: "Hello, ", Span: domain.Span{Start: 0, Length: 9}},
		{Type: domain.UNKNOWN, Value: "a", Span: domain.Span{Start: 9, Length: 1}},
		{Type: interpolatedCloseTt, Value: " world!", Span: domain.Span{Start: 10, Length: 9}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 19, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("\"Hello, {a} world!\""))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i+1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
	}
}

func TestLexMultipleStringParts(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	expected := []domain.Token{
		{Type: interpolatedOpenTt, Value: "Hello, ", Span: domain.Span{Start: 0, Length: 9}},
		{Type: domain.UNKNOWN, Value: "a", Span: domain.Span{Start: 9, Length: 1}},
		{Type: interpolatedMiddleTt, Value: " w", Span: domain.Span{Start: 10, Length: 4}},
		{Type: domain.UNKNOWN, Value: "b", Span: domain.Span{Start: 14, Length: 1}},
		{Type: interpolatedCloseTt, Value: "orld!", Span: domain.Span{Start: 15, Length: 7}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 22, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("\"Hello, {a} w{b}orld!\""))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i+1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
	}
}

func TestLexMultipleInnerBraces(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	expected := []domain.Token{
		{Type: interpolatedOpenTt, Value: "Hello, ", Span: domain.Span{Start: 0, Length: 9}},
		{Type: domain.UNKNOWN, Value: "{", Span: domain.Span{Start: 9, Length: 1}},
		{Type: domain.UNKNOWN, Value: "}", Span: domain.Span{Start: 10, Length: 1}},
		{Type: interpolatedMiddleTt, Value: " w", Span: domain.Span{Start: 11, Length: 4}},
		{Type: domain.UNKNOWN, Value: "{", Span: domain.Span{Start: 15, Length: 1}},
		{Type: domain.UNKNOWN, Value: "{", Span: domain.Span{Start: 16, Length: 1}},
		{Type: domain.UNKNOWN, Value: "}", Span: domain.Span{Start: 17, Length: 1}},
		{Type: domain.UNKNOWN, Value: "}", Span: domain.Span{Start: 18, Length: 1}},
		{Type: interpolatedCloseTt, Value: "orld!", Span: domain.Span{Start: 19, Length: 7}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 26, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("\"Hello, {{}} w{{{}}}orld!\""))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i+1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i+1, "tokens")
	}
}

func TestLexMismatchedOpeningBraces(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	// TODO check if it logs properly
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenizer.NewTokenizer(config, []byte("\"Hello, {a\""))
}

func TestLexMismatchedClosingBraces(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	// TODO check if it logs properly
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenizer.NewTokenizer(config, []byte("\"Hello, a}\""))
}

func TestLexQuoteInsideInterpolatedString(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	// TODO check if it logs properly
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenizer.NewTokenizer(config, []byte("\"Hello\", {a}\""))
}

func TestLexInterpolatedStringWithoutStartingQuote(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	stringTt := config.NewTokenType()
	interpolatedOpenTt := config.NewTokenType()
	interpolatedMiddleTt := config.NewTokenType()
	interpolatedCloseTt := config.NewTokenType()
	stringLiteralMatcher := NewStringLiteralMatcher(stringTt, interpolatedOpenTt, interpolatedMiddleTt, interpolatedCloseTt)

	config.AddMatcher(&stringLiteralMatcher)

	// TODO check if it logs properly
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	tokenizer.NewTokenizer(config, []byte("Hello, {a}\""))
}
