package symbols

import (
	"minimal/minimal-core/built-in/tokenizer"
	"minimal/minimal-core/domain"
	"reflect"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
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

	actual := tokenizer.NewTokenizer(config, []byte("1+2-3"))

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

func TestLexMultiCharSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
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

	actual := tokenizer.NewTokenizer(config, []byte("1+-+-2-+-+3"))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
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

func TestLexUnicodeSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
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

	actual := tokenizer.NewTokenizer(config, []byte("1☘2❤3"))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
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

func TestLexVariationSelector(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
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

	actual := tokenizer.NewTokenizer(config, []byte("1❄️2🔥3"))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
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

func TestLexZeroWidthJoinerSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
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

	actual := tokenizer.NewTokenizer(config, []byte("1🐻‍❄️2🐈‍⬛3"))

	i := 0
	for ; actual.CurrentToken().Type != domain.EOF; i++ {
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
