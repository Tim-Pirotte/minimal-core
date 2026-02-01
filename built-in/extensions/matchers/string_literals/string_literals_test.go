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
