package strings

import (
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messaging"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"testing"
)

type testLexer struct {
    l          *lexer.LexerScheme
    stringType lexer.TokenType
    messenger  *messaging.Messenger
    output     *messaging.TestOutput
}

func getLexer() testLexer {
    l := lexer.NewScheme(1)

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
        l,
        stringType,
    )

    l.AddMatcher(stringMatcher)

    return testLexer{l, stringType, messenger, testOutput}
}

func TestString(t *testing.T) {
    source := `'Hello, World!'`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestEscape(t *testing.T) {
    source := `'Hello,\' World!\\'`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestMultiLine(t *testing.T) {
    source :=
` 'Hello,
     \' World!\\
       '`

    l := getLexer()

    expected := []lexer.Token{
        {Type: lexer.UNKNOWN, Value: source[:1]},
        {Type: l.stringType, Value: source[1:]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestInterpolated(t *testing.T) {
    source := `'H\{{ '\}\'' }W'`

    l := getLexer()

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:5]},
        {Type: lexer.UNKNOWN, Value: source[5:6]},
        {Type: l.stringType, Value: source[6:12]},
        {Type: lexer.UNKNOWN, Value: source[12:13]},
        {Type: l.stringType, Value: source[13:16]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestUnclosedString(t *testing.T) {
    source := `'Hello,`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{
        {
            Message: "String is not terminated with '",
            Severity: messaging.Error,
            Context: []messaging.Span{{Content: source[:1], Note: "The string starts here"}},
            Notes: []string{"The remaining content will be interpreted as the string"},
        },
    })
}

func TestMultipleStrings(t *testing.T) {
    source := `'Hello, World!''Hello, World!'`

    l := getLexer()

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:15]},
        {Type: l.stringType, Value: source[15:]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestMissingClosingBrace(t *testing.T) {
    source := `'{ `

    l := getLexer()

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: l.stringType, Value: source[3:3]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "String interpolation is not terminated with }",
                Severity: messaging.Error,
                Context: []messaging.Span{{Content: source[1:2], Note: "Interpolation starts here"}},
                Notes: []string{"The string will be closed at the end of the source code"},
            },
        },
    )
}

func TestExtraClosingBrace(t *testing.T) {
    source := `'}'`

    l := getLexer()

    expected := []lexer.Token{{Type: l.stringType, Value: source}}

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

type braceSymbolMatcher struct { }

func (b *braceSymbolMatcher) New(_ *lexer.Lexer) lexer.Matcher { return b }

func (*braceSymbolMatcher) Match(l *lexer.Lexer) uint {
    if len(l.Data) - int(l.Position) >= 2 && l.GetNextN(2) == `"}` {
        return 2
    }

    return 0
}

func (*braceSymbolMatcher) Consume(_ *lexer.Lexer, _ uint) {}

func TestDifferentEnclosing(t *testing.T) {
    source := `'{ "} }'`

    l := getLexer()
    l.l.AddMatcher(&braceSymbolMatcher{})

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: lexer.UNKNOWN, Value: source[5:6]},
        {Type: l.stringType, Value: source[6:8]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}

func TestMultipleInterpolations(t *testing.T) {
    source := `'{ } {}'`

    l := getLexer()
    l.l.AddMatcher(&braceSymbolMatcher{})

    expected := []lexer.Token{
        {Type: l.stringType, Value: source[:2]},
        {Type: lexer.UNKNOWN, Value: source[2:3]},
        {Type: l.stringType, Value: source[3:6]},
        {Type: l.stringType, Value: source[6:8]},
    }

    lexer.CheckTokens(t, l.l, expected, source)
    l.messenger.Close()
    l.output.CheckMessages(t, []messaging.Message{})
}
