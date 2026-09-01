package symbols

import (
	"fmt"
	"minimal/minimal-core/built-in/lexer"
)

const byteValueCount = 256

type trieNode struct {
	leaf      bool
	tokenType lexer.TokenType
	children  [byteValueCount]*trieNode
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

func (s *SymbolMatcher) New(_ *lexer.Lexer) lexer.Matcher {
	return s
}

func (s *SymbolMatcher) AddSymbol(l *lexer.LexerScheme, symbol string, tokenType lexer.TokenType) {
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
	node.tokenType = tokenType
}

func (s *SymbolMatcher) Match(l *lexer.Lexer) uint32 {
	pos := uint32(0)
	length := uint32(0)

	char, ok := l.Get(pos)
	node := s.symbols.children[char]

	for ; node != nil; node = node.children[char] {
		pos++

		if node.leaf {
			length = pos
			s.cachedTokenType = lexer.TokenType(node.tokenType)
		}

		if char, ok = l.Get(pos); !ok {
			return length
		}
	}

	return length
}

func (s *SymbolMatcher) Consume(l *lexer.Lexer, length uint32) {
	l.Emit(lexer.Token{Type: s.cachedTokenType, Value: l.GetNextN(length)})
}
