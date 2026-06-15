package symbols

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	tokenizer := lexer.NewLexer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	plus := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'+'",
		DebugName: "Plus",
	})

	symbolMatcher.AddSymbol(tokenizer, "+", plus)

	minus := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(tokenizer, "-", minus)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: plus, Value: "+", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 3, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 4, Length: 1}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "1+2-3")
}

func TestLexMultiCharSymbols(t *testing.T) {
	tokenizer := lexer.NewLexer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	minus := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(tokenizer, "-", minus)

	minusMinus := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'--'",
		DebugName: "MinusMinus",
	})

	symbolMatcher.AddSymbol(tokenizer, "--", minusMinus)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minusMinus, Value: "--", Range: primitives.Range{Start: 3, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 5, Length: 1}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "1-2--3")
}

func TestLexUnicodeSymbols(t *testing.T) {
	tokenizer := lexer.NewLexer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	club := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'☘'",
		DebugName: "Club",
	})

	symbolMatcher.AddSymbol(tokenizer, "☘", club)

	heart := tokenizer.NewTokenType(lexer.TokenTypeMetadata{
		DisplayName: "'❤'",
		DebugName: "Heart",
	})

	symbolMatcher.AddSymbol(tokenizer, "❤", heart)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: club, Value: "☘", Range: primitives.Range{Start: 1, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: heart, Value: "❤", Range: primitives.Range{Start: 5, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 8, Length: 1}},
	}

	lexer.CheckTokens(t, tokenizer, expected, "1☘2❤3")
}
