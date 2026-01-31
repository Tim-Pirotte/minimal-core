package whitespace

import (
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
	"reflect"
	"testing"
)

func TestLexSpace(t *testing.T) {
	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "a", Span: domain.Span{Start: 12, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 26, Length: 0}},
	}

	tokenizerConfig := tokenizer.NewTokenizerConfig()

	wm := NewWhiteSpaceMatcher()

	tokenizerConfig.AddMatcher(&wm)

	actual := tokenizer.NewTokenizer(tokenizerConfig, []byte(" \t \t        a\t            "))

	i := 0
	for ;actual.CurrentToken().Type != domain.EOF; i++ {
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
