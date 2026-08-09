package lexer

import (
	"fmt"
	"io"
	"minimal/minimal-core/built-in/diff"
	"minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/primitives"
	"strconv"
)

type LexerDebugger struct {
    lexer     *Lexer
    output    io.Writer
    messenger *messaging.Messenger
}

// TODO Add syntax highlighting
func NewLexerDebugger(
    lexer *Lexer,
    output io.Writer,
    messenger *messaging.Messenger,
) *LexerDebugger {
    return &LexerDebugger{lexer, output, messenger}
}

func (l *LexerDebugger) DisplayTokens(source string, tokens []Token) {
    for _, token := range tokens {
        if _, err := io.WriteString(l.output, l.StringifyToken(source, token)+"\n"); err != nil {
            l.messenger.Send(
                messaging.Message{
                    Message: "Lexer debugger output write failed",
                    Severity: messaging.Error,
                },
            )

            return
        }
    }
}

func compareTokens(a, b Token) bool {
    return a.Type == b.Type && a.Value == b.Value
}

// Prints a diff of tokens to a writer. Tokens are considered the same if there types and values are equal.
// The range is deliberately ignored since a small change can change all the following ranges.
func (l *LexerDebugger) DisplayTokensDiff(source string, before, after []Token) {
    tokenDiff := diff.GetDiff(before, after, compareTokens)

    for _, diffPart := range tokenDiff {
        prefix := "   "

        switch diffPart.Type {
        case diff.Insert:
            prefix = " + "
        case diff.Delete:
            prefix = " - "
        }

        fmt.Print(prefix)

        if _, err := io.WriteString(l.output, l.StringifyToken(source, diffPart.Value)+"\n"); err != nil {
            l.messenger.Send(messaging.Message{Message: "Lexer debugger output write failed"})

            return
        }
    }
}

func (l *LexerDebugger) StringifyToken(source string, token Token) string {
    tokenTypeMetadata := l.lexer.GetTokenTypeMetadata(token.Type)

    if !primitives.IsSubString(source, token.Value) {
        return fmt.Sprintf(
            "%-20s %-20s     not from source",
            tokenTypeMetadata.DebugName,
            strconv.Quote(token.Value),
        )
    }

    start := uint(primitives.GetStringPtr(token.Value) - primitives.GetStringPtr(source))
    length := uint(len(token.Value))

    return fmt.Sprintf(
        "%-20s %-20s %6d..%-6d (%d)",
        tokenTypeMetadata.DebugName,
        strconv.Quote(token.Value),
        start,
        start+length,
        length,
    )
}
