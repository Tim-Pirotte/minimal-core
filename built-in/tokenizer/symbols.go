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
	var tt domain.TokenType
	l := 0
	ok := false
	
	node := s.symbols

	var pos int
	for pos = 0; node != nil; pos++ {
		i, _ := t.Get(pos)
		node = node.children[i]

		if node != nil && node.leaf {
			tt = node.token
			ok = true
			l = pos + 1
		}
	}

	return ok, uint(l), tt, ""
}
