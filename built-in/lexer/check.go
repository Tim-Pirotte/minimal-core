package lexer

import (
    "fmt"
    "testing"
    "os"
    "unsafe"
    "minimal/minimal-lang/built-in/messenger"
    "minimal/minimal-lang/built-in/outputs/log-renderer"
)

func CheckTokens(t *testing.T, scheme *LexerScheme, expected []Token, text string) {
    actual := make([]Token, 0, len(expected))

    lexer := scheme.Lex(text)

    for current := lexer.Peek(0); current.Type != END; current = lexer.Peek(0) {
        actual = append(actual, current)
        lexer.Advance()
    }

    messenger := messenger.New()
    messenger.AddOutput(logrenderer.New(os.Stdout))

    lexerDebugger := NewDisplayer(scheme, os.Stdout, messenger)

    if len(expected) != len(actual) {
        lexerDebugger.DisplayDiff(text, actual, expected)
        fmt.Println("")
        t.Fatal("Expected", len(expected), "tokens but got", len(actual), "tokens")
    }

    ok := true

    for i := range len(expected) {
        if actual[i].Type != expected[i].Type {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect type)",
            )

            ok = false

            break
        } else if actual[i].Value != expected[i].Value {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect value)",
            )

            ok = false

            break
        } else if unsafe.StringData(actual[i].Value) != unsafe.StringData(expected[i].Value) {
            t.Error(
                "\nExpected\n", lexerDebugger.StringifyToken(text, expected[i]),
                "\nbut got\n", lexerDebugger.StringifyToken(text, actual[i]), "(incorrect string address)",
            )

            ok = false

            break
        }
    }

    if !ok {
        lexerDebugger.DisplayDiff(text, actual, expected)
        fmt.Println()
    }
}
