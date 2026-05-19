package tokenizerv2

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"reflect"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	expected := []Token{{Type: EOF, Value: "", Span: usermessaging.Span{Start: 0, Length: 0}}}

	tokenizerConfig := NewTokenizer()

	actual := tokenizerConfig.Tokenize([]byte(""))

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: EOF, Value: "", Span: usermessaging.Span{Start: 1, Length: 0}},
	}

	tokenizerConfig := NewTokenizer()

	actual := tokenizerConfig.Tokenize([]byte("a"))

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}
