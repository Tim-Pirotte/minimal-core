package tokenizerv2

import (
	"minimal/minimal-core/built-in/primitives"
	"reflect"
	"testing"
)

type TestStopper struct {}

func (t *TestStopper) End(s *TokenizerState) bool {
	return s.Position >= uint(len(s.Data))
}

func TestLexEmpty(t *testing.T) {
	tokenizer := NewTokenizer()
	expected := []Token{}

	actual := tokenizer.Tokenize("", &TestStopper{})

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	tokenizer := NewTokenizer()
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	actual := tokenizer.Tokenize("a", &TestStopper{})

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}
