package list

import (
	"fmt"
	"io"
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
		tokenTypeMetadata, ok := tokenizer.GetTokenTypeMetadata(token.Type)
		name := tokenTypeMetadata.DebugName

		if !ok {
			name = strconv.Itoa(int(token.Type))
		}

		_, err := fmt.Fprintf(
			t.output,
			"%-20s %-20s %6d..%-6d (%d)\n",
			name, 
			token.Value, 
			token.Range.Start, 
			token.Range.Start + token.Range.Length, 
			token.Range.Length,
		)

		if err != nil {
			t.logger.Error().Err(err).Msg("failed writing to output")
		}
	}
}
