package symbols

import (
	"minimal/minimal-core/built-in/tokenizer"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"reflect"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	symbolMatcher := NewSymbolMatcher()

	plus := symbolMatcher.AddSymbol(&config, "+")
	minus := symbolMatcher.AddSymbol(&config, "-")

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: plus, Value: "", Span: usermessaging.Span{Start: 1, Length: 1}},
		{Type: tokenizer.UNKNOWN, Value: "2", Span: usermessaging.Span{Start: 2, Length: 1}},
		{Type: minus, Value: "", Span: usermessaging.Span{Start: 3, Length: 1}},
		{Type: tokenizer.UNKNOWN, Value: "3", Span: usermessaging.Span{Start: 4, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 5, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1+2-3"))

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

func TestLexMultiCharSymbols(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	symbolMatcher := NewSymbolMatcher()

	weirdPlus := symbolMatcher.AddSymbol(&config, "+-+-")
	weirdMinus := symbolMatcher.AddSymbol(&config, "-+-+")

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: weirdPlus, Value: "", Span: usermessaging.Span{Start: 1, Length: 4}},
		{Type: tokenizer.UNKNOWN, Value: "2", Span: usermessaging.Span{Start: 5, Length: 1}},
		{Type: weirdMinus, Value: "", Span: usermessaging.Span{Start: 6, Length: 4}},
		{Type: tokenizer.UNKNOWN, Value: "3", Span: usermessaging.Span{Start: 10, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 11, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1+-+-2-+-+3"))

	i := 0
	for ; actual.CurrentToken().Type != tokenizer.EOF; i++ {
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

	club := symbolMatcher.AddSymbol(&config, "☘")
	heart := symbolMatcher.AddSymbol(&config, "❤")

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: club, Value: "", Span: usermessaging.Span{Start: 1, Length: 3}},
		{Type: tokenizer.UNKNOWN, Value: "2", Span: usermessaging.Span{Start: 4, Length: 1}},
		{Type: heart, Value: "", Span: usermessaging.Span{Start: 5, Length: 3}},
		{Type: tokenizer.UNKNOWN, Value: "3", Span: usermessaging.Span{Start: 8, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 9, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1☘2❤3"))

	i := 0
	for ; actual.CurrentToken().Type != tokenizer.EOF; i++ {
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

	ice := symbolMatcher.AddSymbol(&config, "❄️")
	fire := symbolMatcher.AddSymbol(&config, "🔥")

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: ice, Value: "", Span: usermessaging.Span{Start: 1, Length: 6}},
		{Type: tokenizer.UNKNOWN, Value: "2", Span: usermessaging.Span{Start: 7, Length: 1}},
		{Type: fire, Value: "", Span: usermessaging.Span{Start: 8, Length: 4}},
		{Type: tokenizer.UNKNOWN, Value: "3", Span: usermessaging.Span{Start: 12, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 13, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1❄️2🔥3"))

	i := 0
	for ; actual.CurrentToken().Type != tokenizer.EOF; i++ {
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

	polarBear := symbolMatcher.AddSymbol(&config, "🐻‍❄️")
	blackCat := symbolMatcher.AddSymbol(&config, "🐈‍⬛")

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.UNKNOWN, Value: "1", Span: usermessaging.Span{Start: 0, Length: 1}},
		{Type: polarBear, Value: "", Span: usermessaging.Span{Start: 1, Length: 13}},
		{Type: tokenizer.UNKNOWN, Value: "2", Span: usermessaging.Span{Start: 14, Length: 1}},
		{Type: blackCat, Value: "", Span: usermessaging.Span{Start: 15, Length: 10}},
		{Type: tokenizer.UNKNOWN, Value: "3", Span: usermessaging.Span{Start: 25, Length: 1}},
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 26, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte("1🐻‍❄️2🐈‍⬛3"))

	i := 0
	for ; actual.CurrentToken().Type != tokenizer.EOF; i++ {
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

func TestLexSymbolOutOfBounds(t *testing.T) {
	config := tokenizer.NewTokenizerConfig()
	symbolMatcher := NewSymbolMatcher()

	config.AddMatcher(&symbolMatcher)

	expected := []tokenizer.Token{
		{Type: tokenizer.EOF, Value: "", Span: usermessaging.Span{Start: 0, Length: 0}},
	}

	actual := tokenizer.NewTokenizer(config, []byte(""))

	i := 0
	for ; actual.CurrentToken().Type != tokenizer.EOF; i++ {
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
