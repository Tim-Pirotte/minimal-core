package tokenizer

import (
	"minimal/minimal-core/domain"
	"reflect"
	"testing"
)

func TestLexEmpty(t *testing.T) {
	expected := []domain.Token{{Type: domain.EOF, Value: "", Span: domain.Span{Start: 0, Length: 0}}}

	tokenizerConfig := NewTokenizerConfig([]byte(""))

	actual := tokenizerConfig.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexUnknown(t *testing.T) {
	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "a", Span: domain.Span{Start: 0, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 1, Length: 0}},
	}

	tokenizerConfig := NewTokenizerConfig([]byte("a"))

	actual := tokenizerConfig.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "but got", actual)
	}
}

func TestLexSymbols(t *testing.T) {
	config := NewTokenizerConfig([]byte("1+2-3"))
	symbolMatcher := NewSymbolMatcher()

	plus := symbolMatcher.AddSymbol(&config, "+")
	minus := symbolMatcher.AddSymbol(&config, "-")

	config.AddMatcher(&symbolMatcher)

	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "1", Span: domain.Span{Start: 0, Length: 1}},
		{Type: plus, Value: "", Span: domain.Span{Start: 1, Length: 1}},
		{Type: domain.UNKNOWN, Value: "2", Span: domain.Span{Start: 2, Length: 1}},
		{Type: minus, Value: "", Span: domain.Span{Start: 3, Length: 1}},
		{Type: domain.UNKNOWN, Value: "3", Span: domain.Span{Start: 4, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 5, Length: 0}},
	}

	actual := config.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected\n", expected, "but got\n", actual)
	}
}

func TestLexMultiCharSymbols(t *testing.T) {
	config := NewTokenizerConfig([]byte("1+-+-2-+-+3"))
	symbolMatcher := NewSymbolMatcher()

	weirdPlus := symbolMatcher.AddSymbol(&config, "+-+-")
	weirdMinus := symbolMatcher.AddSymbol(&config, "-+-+")

	config.AddMatcher(&symbolMatcher)

	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "1", Span: domain.Span{Start: 0, Length: 1}},
		{Type: weirdPlus, Value: "", Span: domain.Span{Start: 1, Length: 4}},
		{Type: domain.UNKNOWN, Value: "2", Span: domain.Span{Start: 5, Length: 1}},
		{Type: weirdMinus, Value: "", Span: domain.Span{Start: 6, Length: 4}},
		{Type: domain.UNKNOWN, Value: "3", Span: domain.Span{Start: 10, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 11, Length: 0}},
	}

	actual := config.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected\n", expected, "but got\n", actual)
	}
}

func TestLexUnicodeSymbols(t *testing.T) {
	config := NewTokenizerConfig([]byte("1☘2❤3"))
	symbolMatcher := NewSymbolMatcher()

	plus := symbolMatcher.AddSymbol(&config, "☘")
	minus := symbolMatcher.AddSymbol(&config, "❤")

	config.AddMatcher(&symbolMatcher)

	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "1", Span: domain.Span{Start: 0, Length: 1}},
		{Type: plus, Value: "", Span: domain.Span{Start: 1, Length: 3}},
		{Type: domain.UNKNOWN, Value: "2", Span: domain.Span{Start: 4, Length: 1}},
		{Type: minus, Value: "", Span: domain.Span{Start: 5, Length: 3}},
		{Type: domain.UNKNOWN, Value: "3", Span: domain.Span{Start: 8, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 9, Length: 0}},
	}

	actual := config.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected\n", expected, "but got\n", actual)
	}
}

func TestLexVariationSelector(t *testing.T) {
	config := NewTokenizerConfig([]byte("1❄️2🔥3"))
	symbolMatcher := NewSymbolMatcher()
	
	plus := symbolMatcher.AddSymbol(&config, "❄️")
	minus := symbolMatcher.AddSymbol(&config, "🔥")

	config.AddMatcher(&symbolMatcher)

	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "1", Span: domain.Span{Start: 0, Length: 1}},
		{Type: plus, Value: "", Span: domain.Span{Start: 1, Length: 6}},
		{Type: domain.UNKNOWN, Value: "2", Span: domain.Span{Start: 7, Length: 1}},
		{Type: minus, Value: "", Span: domain.Span{Start: 8, Length: 4}},
		{Type: domain.UNKNOWN, Value: "3", Span: domain.Span{Start: 12, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 13, Length: 0}},
	}

	actual := config.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected\n", expected, "but got\n", actual)
	}
}

func TestLexZeroWidthJoinerSymbols(t *testing.T) {
	config := NewTokenizerConfig([]byte("1🐻‍❄️2🐈‍⬛3"))
	symbolMatcher := NewSymbolMatcher()
	
	plus := symbolMatcher.AddSymbol(&config, "🐻‍❄️")
	minus := symbolMatcher.AddSymbol(&config, "🐈‍⬛")

	config.AddMatcher(&symbolMatcher)

	expected := []domain.Token{
		{Type: domain.UNKNOWN, Value: "1", Span: domain.Span{Start: 0, Length: 1}},
		{Type: plus, Value: "", Span: domain.Span{Start: 1, Length: 13}},
		{Type: domain.UNKNOWN, Value: "2", Span: domain.Span{Start: 14, Length: 1}},
		{Type: minus, Value: "", Span: domain.Span{Start: 15, Length: 10}},
		{Type: domain.UNKNOWN, Value: "3", Span: domain.Span{Start: 25, Length: 1}},
		{Type: domain.EOF, Value: "", Span: domain.Span{Start: 26, Length: 0}},
	}

	actual := config.tokenize()

	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected\n", expected, "but got\n", actual)
	}
}
