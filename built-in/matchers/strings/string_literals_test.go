package strings

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

    stringType := l.NewTokenType(
        lexer.TokenTypeMetadata{DisplayName: "a string literal", DebugName: "String"},
    )

    messenger := messaging.NewMessenger()
    logrenderer := logrendering.NewLogRenderer(os.Stdout)
    logrenderer.Config.RemoveANSI()
    logrenderer.Config.RemoveUnicode()
    messenger.AddOutput(logrenderer)

    testOutput := &messaging.TestOutput{}
    messenger.AddOutput(testOutput)

    stringMatcher := NewStringMatcher(
        messenger,
        stringType,
        []EnclosingSet{
            {openSequence: `"`, closingSequence: `"`},
        },
    )

    l.AddMatcher(stringMatcher)

    return testLexer{l, stringType, messenger, testOutput}
}

func TestString(t *testing.T) {
    source := `"Hello, World!"`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
}

func TestEscape(t *testing.T) {
    source := `"Hello,\" World!\\"`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
}

func TestMultiLine(t *testing.T) {
    source :=
` "Hello,
     \" World!\\
       "`

    l := getLexer()

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: l.stringType, Value: source[1:]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
}

func TestInterpolated(t *testing.T) {
    source := `"Hello,\{{ "\}\"" } World!"`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
}

func TestUnclosedString(t *testing.T) {
    source := `"Hello,`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{
        {
            Reference: "TODO",
            Message: "String is not terminated with a quote",
            Severity: messaging.Error,
            Context: []messaging.Span{{Content: source[:1], Note: "The string starts here"}},
            Notes: []string{"The remaining content will be interpreted as the string"},
        },
    })
}

func TestMultipleStrings(t *testing.T) {
    source := `"Hello, World!""Hello, World!"`

    l := getLexer()

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:15]},
        {Type: l.stringType, Value: source[15:]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
}

func TestMissingClosingBrace(t *testing.T) {

}

func TestDifferentEnclosing(t *testing.T) {
    source := `"{ '}' }"`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}
