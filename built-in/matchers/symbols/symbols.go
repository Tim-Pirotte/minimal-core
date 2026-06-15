package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

type SymbolMatcher struct {
	symbols *trieNode
	cachedTokenType tokenizerv2.TokenType
}

func NewSymbolMatcher() *SymbolMatcher {
	return &SymbolMatcher{
		&trieNode{children: [256]*trieNode{}},
		tokenizerv2.TokenType(0),
	}
}

func (s *SymbolMatcher) AddSymbol(t *tokenizerv2.Tokenizer, symbol string, tokenType tokenizerv2.TokenType) {
	err := updateTrie(s.symbols, symbol, tokenType)

	if err != nil {
		// TODO log error
		fmt.Println(err.Error())
	}
}

func (s *SymbolMatcher) Match(t *tokenizerv2.TokenizerJob) uint {
	pos := uint(0)
	length := uint(0)

	char, ok := t.Get(pos)
	node := s.symbols.children[char]

	for ; node != nil; node = node.children[char] {
		pos++

		if node.leaf {
			length = pos
			s.cachedTokenType = tokenizerv2.TokenType(node.token)
		}

		if char, ok = t.Get(pos); !ok {
			return length
		}
	}

	return length
}

func (s *SymbolMatcher) Consume(t *tokenizerv2.TokenizerJob, length uint) {
	symbol, _ := t.GetRange(t.Position, length)

	t.Emit(tokenizerv2.Token{
		Type: s.cachedTokenType,
		Value: symbol,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
