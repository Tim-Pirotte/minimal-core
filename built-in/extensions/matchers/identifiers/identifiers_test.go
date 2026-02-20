package identifiers

import (
	"minimal/minimal-core/built-in/tokenizer"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"reflect"
	"testing"
)

func TestLexIdentifiers(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	identifierType := config.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	config.AddMatcher(&identifierMatcher)

	expected := []tokenizer.Token{
		{Type: identifierType, Value: "identifier1", Span: usermessaging.Span{Start: 0, Length: 11}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 11, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("identifier1"))

	i := 0
	for ;actual.CurrentToken().Type != tokenizer.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i + 1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
	}
}

func TestLexMultipleIdentifiers(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	identifierType := config.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	config.AddMatcher(&identifierMatcher)

	expected := []tokenizer.Token{
		{Type: identifierType, Value: "identifier1", Span: usermessaging.Span{Start: 0, Length: 11}},
		{Type: tokenizer.UNKNOWN, Value: " ", Span: usermessaging.Span{Start: 11, Length: 1}},
		{Type: identifierType, Value: "identifier2", Span: usermessaging.Span{Start: 12, Length: 11}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 23, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("identifier1 identifier2"))

	i := 0
	for ;actual.CurrentToken().Type != tokenizer.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i + 1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
	}
}

func TestLexUnicode(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	identifierType := config.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	config.AddMatcher(&identifierMatcher)

	expected := []tokenizer.Token{
		{Type: identifierType, Value: "🐻‍❄️", Span: usermessaging.Span{Start: 0, Length: 13}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 13, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("🐻‍❄️"))

	i := 0
	for ;actual.CurrentToken().Type != tokenizer.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i + 1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
	}
}

func TestLexStartingWithNumber(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	identifierType := config.NewTokenType()
	identifierMatcher := NewIdentifierMatcher(identifierType)

	config.AddMatcher(&identifierMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: identifierType, Value: "identifier", Span: usermessaging.Span{Start: 1, Length: 10}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 11, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1identifier"))

	i := 0
	for ;actual.CurrentToken().Type != tokenizer.EOF; i++ {
		if i >= len(expected) {
			t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
		}

		if !reflect.DeepEqual(actual.CurrentToken(), expected[i]) {
			t.Error("Expected", expected[i], "but got", actual.CurrentToken())
		}

		actual.Consume()
	}

	if i + 1 != len(expected) {
		t.Fatal("Expected", len(expected), "tokens but got", i + 1, "tokens")
	}
}
