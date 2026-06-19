package indentation

import (
	"minimal/minimal-core/built-in/lexer"
	eol "minimal/minimal-core/built-in/matchers/end-of-line"
	"minimal/minimal-core/built-in/primitives"
	"testing"
)

func getLexer() (*lexer.Lexer, lexer.TokenType, lexer.TokenType, lexer.TokenType) {
	l := lexer.NewLexer()

	openBlock := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a new block", DebugName: "OpenBlock"},
	)

	closeBlock := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "the end of a block", DebugName: "CloseBlock"},
	)

	indentationMatcher := NewIndentationMatcher(':', ' ', openBlock, closeBlock)
	l.AddMatcher(indentationMatcher)

	eolType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "the end of the line", DebugName: "EOL"},
	)

	eolMatcher := eol.NewEOLMatcher(eolType)
	l.AddMatcher(eolMatcher)

	return l, openBlock, closeBlock, eolType
}

func TestCorrect(t *testing.T) {
	source := `a
b

c:
   1:
      $


   2

d`

	l, openBlock, closeBlock, eolType := getLexer()

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "b", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 3, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "c", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: openBlock, Value: ":", Range: primitives.Range{Start: 5, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 6, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 7, Length: 1}},
		{Type: openBlock, Value: ":", Range: primitives.Range{Start: 8, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 9, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "$", Range: primitives.Range{Start: 10, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: closeBlock, Value: "   ", Range: primitives.Range{Start: 12, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 15, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 16, Length: 1}},
		{Type: closeBlock, Value: "", Range: primitives.Range{Start: 17, Length: 0}},
		{Type: lexer.UNKNOWN, Value: "d", Range: primitives.Range{Start: 17, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}
