package tokenizer

import (
	"fmt"
	"minimal/minimal-core/domain"
)

type SymbolMatcher struct {
	symbols *trieNode
}

func NewSymbolMatcher() SymbolMatcher {
	return SymbolMatcher{&trieNode{children: [256]*trieNode{}}}
}

func (s *SymbolMatcher) AddSymbol(t *TokenizerConfig, symbol string) domain.TokenType {
	tokenType := t.NewTokenType()
	err := updateTrie(s.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
	
	return tokenType
}

func (s *SymbolMatcher) Match(t *TokenizerConfig) (bool, uint, domain.TokenType, string) {
	return false, 0, domain.IGNORE, ""
}
