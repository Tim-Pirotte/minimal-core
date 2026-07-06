package rawstring

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"testing"
)

type testLexer struct {
    l         *lexer.Lexer
    strType   lexer.TokenType
    messenger *messaging.Messenger
    output    *messaging.TestOutput
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

}
