package lexer

import (
	"fmt"
	"io"
	"minimal/minimal-core/built-in/diff"
	logging "minimal/minimal-core/built-in/internal-logging"
	"strconv"
)

type TokenizerDebugger struct {
	tokenizer *Lexer
	logger    logging.Logger
	output    io.Writer
}

func NewTokenizerDebugger(
	tokenizer *Lexer,
	sourceGen *logging.SourceGenerator,
	output io.Writer,
) *TokenizerDebugger {
	logger, _ := sourceGen.GetLogger("tokenListDisplay")

	return &TokenizerDebugger{tokenizer, logger, output}
}

func (t *TokenizerDebugger) DisplayTokens(tokens []Token) {
	for _, token := range tokens {
		if _, err := io.WriteString(t.output, t.StringifyToken(token)+"\n"); err != nil {
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
func (t *TokenizerDebugger) DisplayTokensDiff(before, after []Token) {
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

		if _, err := io.WriteString(t.output, t.StringifyToken(diffPart.Value)+"\n"); err != nil {
			t.logger.Error().Err(err).Msg("failed writing to output")

			return
		}
	}
}

func (t *TokenizerDebugger) StringifyToken(token Token) string {
	tokenTypeMetadata, ok := t.tokenizer.GetTokenTypeMetadata(token.Type)
	name := tokenTypeMetadata.DebugName

	if !ok {
		name = strconv.Itoa(int(token.Type))
	}

	value := token.Value

	if value != "" {
		value = "\"" + value + "\""
	} else {
		value = "  " // To preserve alignment
	}

	return fmt.Sprintf(
		"%-20s %-20s %6d..%-6d (%d)",
		name,
		value,
		token.Range.Start,
		token.Range.Start+token.Range.Length,
		token.Range.Length,
	)
}
