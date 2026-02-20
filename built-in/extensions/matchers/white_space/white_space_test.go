package whitespace

import (
	"minimal/minimal-core/built-in/tokenizer"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"reflect"
	"testing"
)

func TestLexSpace(t *testing.T) {
	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "a", Span: usermessaging.Span{Start: 12, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 26, Length: 0}},
	}

	tokenizerConfig := tokenizer.NewTokenizerConfig()

	wm := NewWhiteSpaceMatcher()

	tokenizerConfig.AddMatcher(&wm)

	actual := tokenizer.NewTokenizer(tokenizerConfig, []byte(" \t \t        a\t            "))

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
