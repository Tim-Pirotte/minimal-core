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

    stringMatcher := NewStringMatcher(messenger, stringType)

    l.AddMatcher(stringMatcher)

    return testLexer{l, stringType, messenger}
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
        {Type: lexer.UNKNOWN, Value: source[0:1]},
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
