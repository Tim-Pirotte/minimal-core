package list

import (
	"bufio"
	"fmt"
	logging "minimal/minimal-core/built-in/internal-logging"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"strconv"
)

type TokenListDisplay struct {
	logger logging.Logger
}

func NewTokenListDisplay(sourceGen logging.SourceGenerator) TokenListDisplay {
	logger, _ := sourceGen.GetLogger("tokenListDisplay")

	return TokenListDisplay{logger}
} 

func (t *TokenListDisplay) displayTokens(tokenizer *tokenizerv2.Tokenizer, tokens []tokenizerv2.Token, output bufio.Writer) {
	for _, token := range tokens {
		tokenTypeMetadata, ok := tokenizer.GetTokenTypeMetadata(token.Type)
		name := tokenTypeMetadata.DebugName

		if !ok {
			name = strconv.Itoa(int(token.Type))
		}

		_, err := fmt.Fprintf(
			&output,
			"%s %s %d..%d (%d)",
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
