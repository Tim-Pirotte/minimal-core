package tokenizer

import (
	"fmt"
	"minimal/minimal-core/domain"
)

const MaxSymbolLengthBytes = 16

type SymbolMatcher struct {
	symbols *trieNode
}

func NewSymbolMatcher() SymbolMatcher {
	return SymbolMatcher{&trieNode{children: [256]*trieNode{}}}
}

func (s *SymbolMatcher) AddSymbol(t *Tokenizer, symbol string) domain.TokenType {
	if len(symbol) > MaxSymbolLengthBytes {
		// TODO log error
	}

	tokenType := t.NewTokenType()
	err := updateTrie(s.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
	
	return tokenType
}
