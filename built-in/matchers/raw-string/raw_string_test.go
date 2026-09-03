package rawstring

import (
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/outputs/log-renderer"
	testoutput "minimal/minimal-lang/built-in/outputs/test"
	"os"
	"testing"
)

type testLexer struct {
    l          *lexer.LexerScheme
    stringType lexer.TokenType
    messenger  *messenger.Messenger
    output     *testoutput.TestOutput
}

func getLexer() testLexer {
    l := lexer.NewScheme()

    rawStringType := l.NewTokenType(
        lexer.TokenTypeMetadata{NounPhrase: "a raw string literal", DebugName: "RawString"},
    )

    m := messenger.New()
    logrenderer := logrenderer.New(os.Stdout)
    logrenderer.Config.RemoveANSI()
    logrenderer.Config.RemoveUnicode()
    m.AddOutput(logrenderer)

    testOutput := testoutput.New()
    m.AddOutput(testOutput)

    stringMatcher := NewRawStringMatcher(
        m,
        rawStringType,
    )

    l.AddMatcher(stringMatcher)

    return testLexer{l, rawStringType, m, testOutput}
}

func TestRawString(t *testing.T) {
    source := `-'Hello,\n World!'-`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestNoEscape(t *testing.T) {
    source := `---'Hello, World!\'---`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{})
}

func TestMissingDash(t *testing.T) {
    source := `---'Hello,'- World!'-- a'-`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()

    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: `Raw string is not terminated with '---`,
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[:4], Note: "The raw string starts here"}},
            Notes: []string{
                "The amount of dashes in the string prefix must match with the suffix",
                "The remaining content will be interpreted as the raw string",
            },
            Suggestions: []messenger.Suggestion{
                {
                    Suggestion: "This looks like an ending sequence",
                    Replacements: []messenger.Replacement{
                        {
                            From: messenger.Span{
                                Content: source[19:22],
                                Note: "Missing 1 dash",
                            },
                            To: messenger.Span{
                                Content: "'---",
                            },
                        },
                    },
                },
            },
        },
    })
}

func TestMissingDashes(t *testing.T) {
    source := `---'Hello, World!'- a`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()

    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: `Raw string is not terminated with '---`,
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[:4], Note: "The raw string starts here"}},
            Notes: []string{
                "The amount of dashes in the string prefix must match with the suffix",
                "The remaining content will be interpreted as the raw string",
            },
            Suggestions: []messenger.Suggestion{
                {
                    Suggestion: "This looks like an ending sequence",
                    Replacements: []messenger.Replacement{
                        {
                            From: messenger.Span{
                                Content: source[17:19],
                                Note: "Missing 2 dashes",
                            },
                            To: messenger.Span{
                                Content: "'---",
                            },
                        },
                    },
                },
            },
        },
    })
}

func TestUnclosed(t *testing.T) {
    source := `---'Hello, World!`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messenger.Message{
        {
            Message: `Raw string is not terminated with '---`,
            Severity: messenger.Error,
            Context: []messenger.Span{{Content: source[:4], Note: "The raw string starts here"}},
            Notes: []string{
                "The amount of dashes in the string prefix must match with the suffix",
                "The remaining content will be interpreted as the raw string",
            },
        },
    })
}
