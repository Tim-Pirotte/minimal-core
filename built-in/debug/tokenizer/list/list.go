package list

import (
	"fmt"
	"io"
	"minimal/minimal-core/built-in/diff"
	logging "minimal/minimal-core/built-in/internal-logging"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"strconv"
)

type TokenListDisplay struct {
	logger logging.Logger
	output io.Writer
}

func NewTokenListDisplay(sourceGen logging.SourceGenerator, output io.Writer) *TokenListDisplay {
	logger, _ := sourceGen.GetLogger("tokenListDisplay")

	return &TokenListDisplay{logger, output}
} 

func (t *TokenListDisplay) DisplayTokens(tokenizer *tokenizerv2.Tokenizer, tokens []tokenizerv2.Token) {
	for _, token := range tokens {
		if !t.displayToken(tokenizer, token) {
			return
		}
	}
}

func compareTokens(a, b tokenizerv2.Token) bool {
	return a.Type == b.Type && a.Value == b.Value
}

// Prints a diff of tokens to stdout. Tokens are considered the same if there types and values are equal.
// The range is deliberately ignored since a small change can change all the following ranges.
func (t *TokenListDisplay) DisplayTokensDiff(tokenizer *tokenizerv2.Tokenizer, before, after []tokenizerv2.Token) {
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

		if !t.displayToken(tokenizer, diffPart.Value) {
			return
		}
	}
}

func (t *TokenListDisplay) displayToken(tokenizer *tokenizerv2.Tokenizer, token tokenizerv2.Token) (ok bool) {
	tokenTypeMetadata, ok := tokenizer.GetTokenTypeMetadata(token.Type)
	name := tokenTypeMetadata.DebugName

	if !ok {
		name = strconv.Itoa(int(token.Type))
	}

	value := token.Value

	if value != "" {
		value = "\"" + value + "\""
	}

	_, err := fmt.Fprintf(
		t.output,
		"%-20s %-20s %6d..%-6d (%d)\n",
		name, 
		value, 
		token.Range.Start, 
		token.Range.Start + token.Range.Length, 
		token.Range.Length,
	)

	if err != nil {
		t.logger.Error().Err(err).Msg("failed writing to output")

		return false
	}

	return true
}
