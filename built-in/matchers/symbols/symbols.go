package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/primitives"
	tokenizerv2 "minimal/minimal-core/built-in/tokenizer-v2"
)

const byteValueCount = 256

type trieNode struct {
	leaf bool
	token tokenizerv2.TokenType
	children [byteValueCount]*trieNode
}

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
	node := s.symbols

	for _, char := range []byte(symbol) {
        if node.children[char] == nil {
            node.children[char] = &trieNode{}
        }

        node = node.children[char]
    }

	if node.leaf {
		// TODO log error message or panic
		fmt.Printf("'%v' has already been declared", symbol)
	}

	node.leaf = true
	node.token = tokenType
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
