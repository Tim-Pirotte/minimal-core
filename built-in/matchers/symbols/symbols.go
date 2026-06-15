package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/primitives"
)

const byteValueCount = 256

type trieNode struct {
	leaf bool
	token lexer.TokenType
	children [byteValueCount]*trieNode
}

type SymbolMatcher struct {
	symbols *trieNode
	cachedTokenType lexer.TokenType
}

func NewSymbolMatcher() *SymbolMatcher {
	return &SymbolMatcher{
		&trieNode{children: [256]*trieNode{}},
		lexer.TokenType(0),
	}
}

func (s *SymbolMatcher) AddSymbol(t *lexer.Lexer, symbol string, tokenType lexer.TokenType) {
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

func (s *SymbolMatcher) Match(t *lexer.LexerJob) uint {
	pos := uint(0)
	length := uint(0)

	char, ok := t.Get(pos)
	node := s.symbols.children[char]

	for ; node != nil; node = node.children[char] {
		pos++

		if node.leaf {
			length = pos
			s.cachedTokenType = lexer.TokenType(node.token)
		}

		if char, ok = t.Get(pos); !ok {
			return length
		}
	}

	return length
}

func (s *SymbolMatcher) Consume(t *lexer.LexerJob, length uint) {
	symbol, _ := t.GetRange(t.Position, length)

	t.Emit(lexer.Token{
		Type: s.cachedTokenType,
		Value: symbol,
		Range: primitives.Range{Start: t.Position, Length: length}},
	)
}
