package lexer

import (
	"fmt"
	"io"
	"minimal/minimal-core/built-in/diff"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/primitives"
	"strconv"
)

type LexerDebugger struct {
    lexer     *Lexer
    logger    logging.Logger
    output    io.Writer
}

func NewLexerDebugger(
    lexer *Lexer,
    sourceGen *logging.SourceGenerator,
    output io.Writer,
) *LexerDebugger {
    logger, _ := sourceGen.GetLogger("TokenListDisplay")

    return &LexerDebugger{lexer, logger, output}
}

func (t *LexerDebugger) DisplayTokens(source string, tokens []Token) {
    for _, token := range tokens {
        if _, err := io.WriteString(t.output, t.StringifyToken(source, token)+"\n"); err != nil {
            t.logger.Error().Err(err).Msg("failed writing to output")

            return
        }
    }
}

func compareTokens(a, b Token) bool {
    return a.Type == b.Type && a.Value == b.Value
}

// Prints a diff of tokens to a writer. Tokens are considered the same if there types and values are equal.
// The range is deliberately ignored since a small change can change all the following ranges.
func (t *LexerDebugger) DisplayTokensDiff(source string, before, after []Token) {
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

        if _, err := io.WriteString(t.output, t.StringifyToken(source, diffPart.Value)+"\n"); err != nil {
            t.logger.Error().Err(err).Msg("failed writing to output")

            return
        }
    }
}

func (t *LexerDebugger) StringifyToken(source string, token Token) string {
    tokenTypeMetadata, ok := t.lexer.GetTokenTypeMetadata(token.Type)
    name := tokenTypeMetadata.DebugName

    if !ok {
        name = strconv.Itoa(int(token.Type))
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
