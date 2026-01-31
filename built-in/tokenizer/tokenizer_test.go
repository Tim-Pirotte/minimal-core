package tokenizer

import (
	"minimal/minimal-core/domain"
	"reflect"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	expected := []domain.Token{{Type: domain.EOF, Value: "", Span: domain.Span{Start: 0, Length: 0}}}

	tokenizerConfig := NewTokenizerConfig()

	actual := tokenizerConfig.tokenize([]byte(""))

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "a", Span: domain.Span{Start: 0, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 1, Length: 0}},
	}

	tokenizerConfig := NewTokenizerConfig()

	actual := tokenizerConfig.tokenize([]byte("a"))

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}
