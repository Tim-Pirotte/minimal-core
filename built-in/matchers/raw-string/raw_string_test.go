package rawstring

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"testing"
)

type testLexer struct {
    l          *lexer.Lexer
    stringType lexer.TokenType
    messenger  *messaging.Messenger
    output     *messaging.TestOutput
}

func getLexer() testLexer {
    l := lexer.NewLexer()

    rawStringType := l.NewTokenType(
        lexer.TokenTypeMetadata{DisplayName: "a raw string literal", DebugName: "RawString"},
    )

    messenger := messaging.NewMessenger()
    logrenderer := logrendering.NewLogRenderer(os.Stdout)
    logrenderer.Config.RemoveANSI()
    logrenderer.Config.RemoveUnicode()
    messenger.AddOutput(logrenderer)

    testOutput := &messaging.TestOutput{}
    messenger.AddOutput(testOutput)

    stringMatcher := NewRawStringMatcher(
        messenger,
        rawStringType,
    )

    l.AddMatcher(stringMatcher)

    return testLexer{l, rawStringType, messenger, testOutput}
}

func TestRawString(t *testing.T) {
    source := `-"Hello,\n World!"-`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestNoEscape(t *testing.T) {
    source := `---"Hello, World!\"---`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

// TODO improve error message to catch this specific case
func TestMissingDash(t *testing.T) {
    source := `---"Hello, World!"-- a`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestUnclosed(t *testing.T) {
    source := `---"Hello, World!`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{
        {
            Reference: "TODO",
            Message: `Raw string is not terminated with the "--- sequence`,
            Severity: messaging.Error,
            Context: []messaging.Span{{Content: source[:4], Note: "The raw string starts here"}},
            Notes: []string{
                "The amount of dashes in the string prefix must match with the suffix",
                "The remaining content will be interpreted as the raw string",
            },
        },
    })
}
