package symbols

import (
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"testing"
)

func TestLexSymbols(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	plus := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'+'",
		DebugName: "Plus",
	})

	symbolMatcher.AddSymbol(tokenizer, "+", plus)

	minus := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(tokenizer, "-", minus)

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: plus, Value: "+", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: tokenizerv2.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 3, Length: 1}},
		{Type: tokenizerv2.UNKNOWN, Value: "3", Range: primitives.Range{Start: 4, Length: 1}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "1+2-3")
}

func TestLexMultiCharSymbols(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	minus := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'-'",
		DebugName: "Minus",
	})

	symbolMatcher.AddSymbol(tokenizer, "-", minus)

	minusMinus := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'--'",
		DebugName: "MinusMinus",
	})

	symbolMatcher.AddSymbol(tokenizer, "--", minusMinus)

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: minus, Value: "-", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: tokenizerv2.UNKNOWN, Value: "2", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: minusMinus, Value: "--", Range: primitives.Range{Start: 3, Length: 2}},
		{Type: tokenizerv2.UNKNOWN, Value: "3", Range: primitives.Range{Start: 5, Length: 1}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "1-2--3")
}

func TestLexUnicodeSymbols(t *testing.T) {
	tokenizer := tokenizerv2.NewTokenizer()
	symbolMatcher := NewSymbolMatcher()
	tokenizer.AddMatcher(symbolMatcher)

	club := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'☘'",
		DebugName: "Club",
	})

	symbolMatcher.AddSymbol(tokenizer, "☘", club)

	heart := tokenizer.NewTokenType(tokenizerv2.TokenTypeMetadata{
		DisplayName: "'❤'",
		DebugName: "Heart",
	})

	symbolMatcher.AddSymbol(tokenizer, "❤", heart)

	expected := []tokenizerv2.Token{
		{Type: tokenizerv2.UNKNOWN, Value: "1", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: club, Value: "☘", Range: primitives.Range{Start: 1, Length: 3}},
		{Type: tokenizerv2.UNKNOWN, Value: "2", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: heart, Value: "❤", Range: primitives.Range{Start: 5, Length: 3}},
		{Type: tokenizerv2.UNKNOWN, Value: "3", Range: primitives.Range{Start: 8, Length: 1}},
	}

	tokenizerv2.CheckTokens(t, tokenizer, expected, "1☘2❤3")
}
