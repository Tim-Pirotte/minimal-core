package symbols

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	l := lexer.NewLexer()
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
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: plus, Value: "+", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 3, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 4, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, "1+2-3")
}

func TestLexMultiCharSymbols(t *testing.T) {
	l := lexer.NewLexer()
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
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minusMinus, Value: "--", Range: primitives.Range{Start: 3, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 5, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, "1-2--3")
}

func TestLexUnicodeSymbols(t *testing.T) {
	l := lexer.NewLexer()
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
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: club, Value: "☘", Range: primitives.Range{Start: 1, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: heart, Value: "❤", Range: primitives.Range{Start: 5, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 8, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, "1☘2❤3")
}
