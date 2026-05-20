package tokenizerv2

import (
	"minimal/minimal-core/built-in/primitives"
	"reflect"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	expected := []Token{{Type: EOF, Value: "", Range: primitives.Range{Start: 0, Length: 0}}}

	tokenizerConfig := NewTokenizer()

	actual := tokenizerConfig.Tokenize("")

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: EOF, Value: "", Range: primitives.Range{Start: 1, Length: 0}},
	}

	tokenizerConfig := NewTokenizer()

	actual := tokenizerConfig.Tokenize("a")

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}
