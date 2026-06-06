package tokenizerv2

import (
	"minimal/minimal-core/built-in/primitives"
	"reflect"
	"testing"
)

type testStopper struct {}

func (t *testStopper) End(s *TokenizerJob) bool {
	return s.Position >= uint(len(s.Data))
}

func TestLexEmpty(t *testing.T) {
	tokenizer := NewTokenizer()
	expected := []Token{}

	actual := tokenizer.Tokenize("", &testStopper{}, 1)

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	tokenizer := NewTokenizer()
	expected := []Token{
		{Type: UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	tokenizerJob := tokenizer.Tokenize("a", &testStopper{}, 1)

	actual := []Token{}

	for current := tokenizerJob.Peek(0); current.Type != END; current = tokenizerJob.Peek(0) {
		actual = append(actual, current)
		tokenizerJob.Advance()
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func Benchmark(b *testing.B) {
	source := string(make([]byte, 1_000_000))
	tokenizer := NewTokenizer()

	for b.Loop() {
		tokenizer.Tokenize(source, &testStopper{}, 1)
	}
}
