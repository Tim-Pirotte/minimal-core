package list

import (
	"bufio"
	"fmt"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
	"strconv"
)

func displayTokens(tokenizer *tokenizerv2.Tokenizer, tokens []tokenizerv2.Token, output bufio.Writer) {
	for _, token := range tokens {
		tokenTypeMetadata, ok := tokenizer.GetTokenTypeMetadata(token.Type)
		name := tokenTypeMetadata.DebugName

		if !ok {
			name = strconv.Itoa(int(token.Type))
		}

		fmt.Fprintf(
			&output,
			"%s %s %d..%d (%d)",
			name, 
			token.Value, 
			token.Range.Start, 
			token.Range.Start + token.Range.Length, 
			token.Range.Length,
		)
	}
}
