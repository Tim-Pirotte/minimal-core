package lexer

import (
	"fmt"
	"io"
	"minimal/minimal-core/built-in/ansi"
	"minimal/minimal-core/built-in/diff"
	"minimal/minimal-core/built-in/messaging"
	"minimal/minimal-core/built-in/primitives"
	"strconv"
)

type LexerDisplayer struct {
    lexer       *LexerScheme
    output      io.Writer
    messenger   *messaging.Messenger
    tokenColors map[TokenType] ansi.RGB
}

func NewLexerDisplayer(
    scheme *LexerScheme,
    output io.Writer,
    messenger *messaging.Messenger,
) *LexerDisplayer {
    return &LexerDisplayer{scheme, output, messenger, map[TokenType]ansi.RGB{}}
}

func (l *LexerDisplayer) SetTokenTypeColor(tokenType TokenType, color ansi.RGB) {
    l.tokenColors[tokenType] = color
}

func (l *LexerDisplayer) DisplayTokens(source string, tokens []Token) {
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
func (l *LexerDisplayer) DisplayTokensDiff(source string, before, after []Token) {
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

func (l *LexerDisplayer) StringifyToken(source string, token Token) string {
    name := l.lexer.GetTokenTypeMetadata(token.Type).DebugName

    if color, ok := l.tokenColors[token.Type]; ok {
        name = string(color) + name + ansi.Reset
    }

    if !primitives.IsSubString(source, token.Value) {
        return fmt.Sprintf(
            "%-20s %-20s     not from source",
            name,
            strconv.Quote(token.Value),
        )
    }

    start := uint(primitives.GetStringPtr(token.Value) - primitives.GetStringPtr(source))
    length := uint(len(token.Value))

    return fmt.Sprintf(
        "%-20s %-20s %6d..%-6d (%d)",
        name,
        strconv.Quote(token.Value),
        start,
        start+length,
        length,
    )
}
