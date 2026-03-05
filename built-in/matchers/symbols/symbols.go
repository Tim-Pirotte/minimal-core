package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/tokenizer"
)

type SymbolMatcher struct {
	symbols *trieNode
}

func NewSymbolMatcher() SymbolMatcher {
	return SymbolMatcher{&trieNode{children: [256]*trieNode{}}}
}

func (s *SymbolMatcher) AddSymbol(t *tokenizer.TokenizerConfig, symbol string) tokenizer.TokenType {
	tokenType := t.NewTokenType()
	err := updateTrie(s.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
	
	return tokenType
}

func (s *SymbolMatcher) Match(so *tokenizer.Source) (uint, tokenizer.TokenType, string) {
	var tt tokenizer.TokenType
	l := 0	
	node := s.symbols

	var pos int
	for pos = 0; node != nil; pos++ {
		i, ok := so.Get(pos)

		if !ok {
			break
		}

		node = node.children[i]

		if node != nil && node.leaf {
			tt = node.token
			l = pos + 1
		}
	}

	return uint(l), tt, ""
}
