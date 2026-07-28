package symbols

import (
	"minimal/minimal-core/built-in/lexer"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	source := "1+2-3"

	l := lexer.New()
	symbolMatcher := NewSymbolMatcher()
	l.AddMatcher(symbolMatcher)

	plus := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'+'",
		DebugName: "Plus",
	})

	symbolMatcher.AddSymbol(l, "+", plus)

	minus := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(l, "-", minus)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: plus, Value: source[1:2]},
		{Type: lexer.UNKNOWN, Value: source[2:3]},
		{Type: minus, Value: source[3:4]},
		{Type: lexer.UNKNOWN, Value: source[4:5]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestLexMultiCharSymbols(t *testing.T) {
	source := "1-2--3"

	l := lexer.New()
	symbolMatcher := NewSymbolMatcher()
	l.AddMatcher(symbolMatcher)

	minus := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(l, "-", minus)

	minusMinus := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'--'",
		DebugName: "MinusMinus",
	})

	symbolMatcher.AddSymbol(l, "--", minusMinus)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: minus, Value: source[1:2]},
		{Type: lexer.UNKNOWN, Value: source[2:3]},
		{Type: minusMinus, Value: source[3:5]},
		{Type: lexer.UNKNOWN, Value: source[5:6]},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestLexUnicodeSymbols(t *testing.T) {
	source := "1☘2❤3"

	l := lexer.New()
	symbolMatcher := NewSymbolMatcher()
	l.AddMatcher(symbolMatcher)

	club := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'☘'",
		DebugName: "Club",
	})

	symbolMatcher.AddSymbol(l, "☘", club)

	heart := l.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'❤'",
		DebugName: "Heart",
	})

	symbolMatcher.AddSymbol(l, "❤", heart)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: source[:1]},
		{Type: club, Value: source[1:4]},
		{Type: lexer.UNKNOWN, Value: source[4:5]},
		{Type: heart, Value: source[5:8]},
		{Type: lexer.UNKNOWN, Value: source[8:9]},
	}

	lexer.CheckTokens(t, l, expected, source)
}
