package indentation

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messenger"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	testoutput "minimal/minimal-core/built-in/outputs/test"
	"os"
	"testing"
)

type testLexer struct {
    l          *lexer.LexerScheme
    openBlock  lexer.TokenType
    closeBlock lexer.TokenType
    eolType    lexer.TokenType
    messenger  *messenger.Messenger
    output     *testoutput.TestOutput
}

func getLexer(indentChar byte, spacesPerLevel uint) testLexer {
    l := lexer.NewScheme()

    openBlock := l.NewTokenType(
        lexer.TokenTypeMetadata{DisplayName: "a new block", DebugName: "OpenBlock"},
    )

    closeBlock := l.NewTokenType(
        lexer.TokenTypeMetadata{DisplayName: "the end of a block", DebugName: "CloseBlock"},
    )

    eolType := l.NewTokenType(
        lexer.TokenTypeMetadata{DisplayName: "the end of the line", DebugName: "EOL"},
    )

    m := messenger.New()
    logrenderer := logrendering.New(os.Stdout)
    logrenderer.Config.RemoveANSI()
    logrenderer.Config.RemoveUnicode()
    m.AddOutput(logrenderer)

    testOutput := testoutput.New()
    m.AddOutput(testOutput)

    indentationMatcher := NewIndentationMatcher(
        m,
        ':',
        indentChar,
        openBlock,
        closeBlock,
        eolType,
        spacesPerLevel,
    )

    l.AddMatcher(indentationMatcher)

    return testLexer{l, openBlock, closeBlock, eolType, m, testOutput}
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

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: l.eolType, Value: source[1:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: l.eolType, Value: source[3:5]},
        {Type: lexer.UNKNOWN, Value: source[5:6]},
        {Type: l.openBlock, Value: source[6:11]},
        {Type: lexer.UNKNOWN, Value: source[11:12]},
        {Type: l.openBlock, Value: source[12:20]},
        {Type: lexer.UNKNOWN, Value: source[20:21]},
        {Type: l.closeBlock, Value: source[21:27]},
        {Type: lexer.UNKNOWN, Value: source[27:28]},
        {Type: l.closeBlock, Value: source[28:30]},
        {Type: lexer.UNKNOWN, Value: source[30:31]},
        {Type: l.openBlock, Value: source[31:33]},
        {Type: l.closeBlock, Value: source[31:33]},
        {Type: lexer.UNKNOWN, Value: source[33:34]},
        {Type: l.eolType, Value: source[34:36]},
        {Type: lexer.UNKNOWN, Value: source[36:37]},
        {Type: l.openBlock, Value: source[37:42]},
        {Type: lexer.UNKNOWN, Value: source[42:43]},
        {Type: l.openBlock, Value: source[43:51]},
        {Type: lexer.UNKNOWN, Value: source[51:52]},
        {Type: l.closeBlock, Value: source[52:54]},
        {Type: l.closeBlock, Value: source[52:54]},
        {Type: lexer.UNKNOWN, Value: source[54:55]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestRedundantSpace(t *testing.T) {
    source := ":\n"+
              "  a\n" +
              "   \n" +
              "  b"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:4]},
        {Type: lexer.UNKNOWN, Value: source[4:5]},
        {Type: l.eolType, Value: source[5:12]},
        {Type: lexer.UNKNOWN, Value: source[12:13]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestInconsistentIndentation(t *testing.T) {
    source := ":\n"+
              "  :\n" +
              "     a\n" +
              " b"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:4]},
        {Type: l.openBlock, Value: source[4:11]},
        {Type: lexer.UNKNOWN, Value: source[11:12]},
        {Type: l.eolType, Value: source[12:14]},
        {Type: lexer.UNKNOWN, Value: source[14:15]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: "Indentation is inconsistent",
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[6:11]}},
            AdditionalContext: []messenger.Span{
                {
                    Content: source[2:4],
                    Note: "The indentation was derived here",
                },
            },
            Notes: []string{
                "The indentation must be a multiple of 2",
                "The indentation of the incorrect line is 5",
                "The indentation will be set to the current level",
            },
        },
        {
            Message: "Indentation is inconsistent",
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[13:14]}},
            AdditionalContext: []messenger.Span{
                {
                    Content: source[2:4],
                    Note: "The indentation was derived here",
                },
            },
            Notes: []string{
                "The indentation must be a multiple of 2",
                "The indentation of the incorrect line is 1",
                "The indentation will be set to the current level",
            },
        },
    })
}

func TestExtraIndent(t *testing.T) {
    source := ":\n"+
              "  (\n" +
              "     a\n" +
              "  )"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:4]},
        {Type: lexer.UNKNOWN, Value: source[4:5]},
        {Type: l.eolType, Value: source[5:11]},
        {Type: lexer.UNKNOWN, Value: source[11:12]},
        {Type: l.eolType, Value: source[12:15]},
        {Type: lexer.UNKNOWN, Value: source[15:16]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestTooMuchIndent(t *testing.T) {
    source := ":\n"+
              "  :\n" +
              "      a"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:4]},
        {Type: l.openBlock, Value: source[4:12]},
        {Type: lexer.UNKNOWN, Value: source[12:13]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: "Got more indentation than expected",
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[6:12]}},
            AdditionalContext: []messenger.Span{},
            Notes: []string{
                "The indentation of the incorrect line is 6",
                "The largest expected indentation is 4",
                "The indentation will be set to the current level",
            },
        },
    })
}

func TestDifferentSpaceChar(t *testing.T) {
    source := ":\n"+
              "\t:\n" +
              "\t\ta"

    l := getLexer('\t', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:3]},
        {Type: l.openBlock, Value: source[3:7]},
        {Type: lexer.UNKNOWN, Value: source[7:8]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestFixedIndentation(t *testing.T) {
    source := ":\n"+
              "  a"

    l := getLexer(' ', 4)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source[:4]},
        {Type: lexer.UNKNOWN, Value: source[4:5]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: "Indentation is inconsistent",
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[2:4]}},
            AdditionalContext: []messenger.Span{},
            Notes: []string{
                "The indentation was manually set",
                "The indentation must be a multiple of 4",
                "The indentation of the incorrect line is 2",
                "The indentation will be set to the current level",
            },
        },
    })
}

func TestMatchNonIndentSpace(t *testing.T) {
    source := "a b"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: lexer.UNKNOWN, Value: source[1:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestNoBlocks(t *testing.T) {
    source := "a\n" +
              " b"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: l.eolType, Value: source[1:3]},
        {Type: lexer.UNKNOWN, Value: source[3:4]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestNoBlocksIncorrect(t *testing.T) {
    source := " a"

    l := getLexer(' ', 0)

    expected := []lexer.Token{{Type: lexer.UNKNOWN, Value: source[1:]}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: "Source code cannot start with indentation",
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[:1]}},
            Notes: []string{"The indentation at the start will be skipped"},
        },
    })
}

func TestTrailingSpace(t *testing.T) {
    source := "a "

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: lexer.UNKNOWN, Value: source[1:2]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestTrailingEOL(t *testing.T) {
    source := "a\n"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: l.eolType, Value: source[1:2]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestOnlyOpenBlock(t *testing.T) {
    source := ":"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source},
        {Type: l.closeBlock, Value: source},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestOpenBlockSpaceSuffix(t *testing.T) {
    source := ": "

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source},
        {Type: l.closeBlock, Value: source},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestOpenBlockEOLSuffix(t *testing.T) {
    source := ":\n"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.openBlock, Value: source},
        {Type: l.closeBlock, Value: source},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestEOLPrefix(t *testing.T) {
    source := "\na"

    l := getLexer(' ', 0)

    expected := []lexer.Token{
        {Type: l.eolType, Value: source[:1]},
        {Type: lexer.UNKNOWN, Value: source[1:2]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func FuzzIndent(f *testing.F) {
    f.Add("")
    f.Add("\n")
    f.Add(" ")
    f.Add(":\n a")
    f.Add(":\n :\n  a\nb")
    f.Add("a\n b")

    l := getLexer(' ', 0)

    f.Fuzz(func(t *testing.T, source string) {
        l.l.Lex(source)
    })
}
