package indentation

import (
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/primitives"
	eofstopper "minimal/minimal-core/built-in/stoppers/eof-stopper"
	"os"
	"testing"
)

func getLexer(
	indentChar byte,
	spacesPerLevel uint,
) (*lexer.Lexer, lexer.TokenType, lexer.TokenType, lexer.TokenType) {
	l := lexer.NewLexer()

	openBlock := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "a new block", DebugName: "OpenBlock"},
	)

	closeBlock := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "the end of a block", DebugName: "CloseBlock"},
	)

	eolType := l.NewTokenType(
		lexer.TokenTypeMetadata{DisplayName: "the end of the line", DebugName: "EOL"},
	)

	messenger := messaging.NewMessenger(logging.GetTestLogSource(io.Discard))
	messenger.AddOutput(logrendering.NewLogRenderer(logging.GetTestLogSource(io.Discard), os.Stdout))

	indentationMatcher := NewIndentationMatcher(
		messenger,
		':',
		indentChar,
		openBlock,
		closeBlock,
		eolType,
		spacesPerLevel,
	)

	l.AddMatcher(indentationMatcher)

	return l, openBlock, closeBlock, eolType
}

func TestIndentation(t *testing.T) {
	source := `a
b

c:
   1:
      $


   2

d:
e

1:
   2:
      3

1`

	l, openBlock, closeBlock, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "b", Range: primitives.Range{Start: 2, Length: 1}},
		{Type: eolType, Value: "\n\n", Range: primitives.Range{Start: 3, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "c", Range: primitives.Range{Start: 5, Length: 1}},
		{Type: openBlock, Value: ":\n   ", Range: primitives.Range{Start: 6, Length: 5}},
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: openBlock, Value: ":\n      ", Range: primitives.Range{Start: 12, Length: 8}},
		{Type: lexer.UNKNOWN, Value: "$", Range: primitives.Range{Start: 20, Length: 1}},
		{Type: closeBlock, Value: "\n\n\n   ", Range: primitives.Range{Start: 21, Length: 6}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 27, Length: 1}},
		{Type: closeBlock, Value: "\n\n", Range: primitives.Range{Start: 28, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "d", Range: primitives.Range{Start: 30, Length: 1}},
		{Type: openBlock, Value: ":\n", Range: primitives.Range{Start: 31, Length: 2}},
		{Type: closeBlock, Value: ":\n", Range: primitives.Range{Start: 31, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "e", Range: primitives.Range{Start: 33, Length: 1}},
		{Type: eolType, Value: "\n\n", Range: primitives.Range{Start: 34, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 36, Length: 1}},
		{Type: openBlock, Value: ":\n   ", Range: primitives.Range{Start: 37, Length: 5}},
		{Type: lexer.UNKNOWN, Value: "2", Range: primitives.Range{Start: 42, Length: 1}},
		{Type: openBlock, Value: ":\n      ", Range: primitives.Range{Start: 43, Length: 8}},
		{Type: lexer.UNKNOWN, Value: "3", Range: primitives.Range{Start: 51, Length: 1}},
		{Type: closeBlock, Value: "\n\n", Range: primitives.Range{Start: 52, Length: 2}},
		{Type: closeBlock, Value: "\n\n", Range: primitives.Range{Start: 52, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "1", Range: primitives.Range{Start: 54, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestRedundantSpace(t *testing.T) {
	source := ":\n"+
              "  a\n" +
              "   \n" +
              "  b"

	l, openBlock, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n  ", Range: primitives.Range{Start: 0, Length: 4}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: eolType, Value: "\n   \n  ", Range: primitives.Range{Start: 5, Length: 7}},
		{Type: lexer.UNKNOWN, Value: "b", Range: primitives.Range{Start: 12, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestInconsistentIndentation(t *testing.T) {
	source := ":\n"+
              "  :\n" +
              "     a\n" +
			  " b"

	l, openBlock, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n  ", Range: primitives.Range{Start: 0, Length: 4}},
		{Type: lexer.UNKNOWN, Value: "(", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: eolType, Value: "\n     ", Range: primitives.Range{Start: 5, Length: 6}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: eolType, Value: "\n  ", Range: primitives.Range{Start: 12, Length: 3}},
		{Type: lexer.UNKNOWN, Value: ")", Range: primitives.Range{Start: 15, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestExtraIndent(t *testing.T) {
	source := ":\n"+
              "  (\n" +
              "     a\n" +
              "  )"

	l, openBlock, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n  ", Range: primitives.Range{Start: 0, Length: 4}},
		{Type: lexer.UNKNOWN, Value: "(", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: eolType, Value: "\n     ", Range: primitives.Range{Start: 5, Length: 6}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 11, Length: 1}},
		{Type: eolType, Value: "\n  ", Range: primitives.Range{Start: 12, Length: 3}},
		{Type: lexer.UNKNOWN, Value: ")", Range: primitives.Range{Start: 15, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestExtraIndentAfterBlock(t *testing.T) {
	// TODO check the errors
	source := ":\n"+
              "  :\n" +
              "    a\n" +
              "   )"

	l, openBlock, _, eolType := getLexer(' ', 0)

	// TODO expect incorrect indentation
	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n  ", Range: primitives.Range{Start: 0, Length: 4}},
		{Type: lexer.UNKNOWN, Value: "(", Range: primitives.Range{Start: 4, Length: 1}},
		{Type: eolType, Value: "\n    ", Range: primitives.Range{Start: 5, Length: 5}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 10, Length: 1}},
		{Type: eolType, Value: "\n   ", Range: primitives.Range{Start: 11, Length: 4}},
		{Type: lexer.UNKNOWN, Value: ")", Range: primitives.Range{Start: 15, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestDifferentSpaceChar(t *testing.T) {
	source := ":\n"+
              "\t:\n" +
              "\t\ta"

	l, openBlock, _, _ := getLexer('\t', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n\t", Range: primitives.Range{Start: 0, Length: 3}},
		{Type: openBlock, Value: ":\n\t\t", Range: primitives.Range{Start: 3, Length: 4}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 7, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestFixedIndentation(t *testing.T) {
	source := ":\n"+
              "  a"

	l, openBlock, _, _ := getLexer(' ', 4)

	// TODO expect incorrect indentation
	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n  ", Range: primitives.Range{Start: 0, Length: 3}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 7, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestMatchNonIndentSpace(t *testing.T) {
	source := "a b"

	l, _, _, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: lexer.UNKNOWN, Value: " ", Range: primitives.Range{Start: 1, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "b", Range: primitives.Range{Start: 2, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestNoBlocks(t *testing.T) {
	source := "a\n" +
			  " b"

	l, _, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: eolType, Value: "\n ", Range: primitives.Range{Start: 1, Length: 2}},
		{Type: lexer.UNKNOWN, Value: "b", Range: primitives.Range{Start: 3, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestNoBlocksIncorrect(t *testing.T) {
	source := " a"

	l, _, _, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestTrailingSpace(t *testing.T) {
	source := "a "

	l, _, _, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: lexer.UNKNOWN, Value: " ", Range: primitives.Range{Start: 1, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestTrailingEOL(t *testing.T) {
	source := "a\n"

	l, _, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 1, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestOnlyOpenBlock(t *testing.T) {
	source := ":"

	l, openBlock, closeBlock, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: closeBlock, Value: ":", Range: primitives.Range{Start: 0, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestOpenBlockSpaceSuffix(t *testing.T) {
	source := ": "

	l, openBlock, closeBlock, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ": ", Range: primitives.Range{Start: 0, Length: 2}},
		{Type: closeBlock, Value: ": ", Range: primitives.Range{Start: 0, Length: 2}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestOpenBlockEOLSuffix(t *testing.T) {
	source := ":\n"

	l, openBlock, closeBlock, _ := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: openBlock, Value: ":\n", Range: primitives.Range{Start: 0, Length: 2}},
		{Type: closeBlock, Value: ":\n", Range: primitives.Range{Start: 0, Length: 2}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func TestEOLPrefix(t *testing.T) {
	source := "\na"

	l, _, _, eolType := getLexer(' ', 0)

	expected := []lexer.Token{
		{Type: eolType, Value: "\n", Range: primitives.Range{Start: 0, Length: 1}},
		{Type: lexer.UNKNOWN, Value: "a", Range: primitives.Range{Start: 1, Length: 1}},
	}

	lexer.CheckTokens(t, l, expected, source)
}

func FuzzIndent(f *testing.F) {
	f.Add("")
	f.Add("\n")
	f.Add(" ")
	f.Add(":\n a")
	f.Add(":\n :\n  a\nb")
	f.Add("a\n b")

	l, _, _, _ := getLexer(' ', 0)

	f.Fuzz(func(t *testing.T, source string) {
		l.Lex(source, eofstopper.NewEOFStopper(), 1)
	})
}
