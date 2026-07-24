package prefix

import (
	"fmt"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
)

type Rule struct {
    TokenTypes []lexer.TokenType
    Handler    func(*lexer.LexerJob, *ast.AST)
}

type trieNode struct {
    leaf     bool
    handler  func(*lexer.LexerJob, *ast.AST)
    children map[lexer.TokenType]*trieNode
}

type PrefixParser struct {
    prefixes  *trieNode
    maxLength uint
}

func NewPrefixParser(prefixes []Rule) *PrefixParser {
    root := &trieNode{false, nil, map[lexer.TokenType]*trieNode{}}
    maxLength := uint(0)

    for _, prefix := range prefixes {
        length := uint(len(prefix.TokenTypes))

        if length > maxLength {
            maxLength = length
        }

        node := root

        for _, tokenType := range prefix.TokenTypes {
            if _, ok := node.children[tokenType]; !ok {
                node.children[tokenType] = &trieNode{children: map[lexer.TokenType]*trieNode{}}
            }

            node = node.children[tokenType]
        }

        if node.leaf {
            // TODO log error message or panic
            panic(fmt.Sprintf("'%v' has already been declared", prefix))
        }

        node.leaf = true
        node.handler = prefix.Handler
    }

    return &PrefixParser{root, uint(maxLength)}
}

func (p *PrefixParser) Parse(lj *lexer.LexerJob, syntax *ast.AST) {
    var largestMatchHandler func(*lexer.LexerJob, *ast.AST)
    node := p.prefixes
    ok := true

    if node.leaf {
        largestMatchHandler = node.handler
    }

    for pos := uint(0); ok && pos < p.maxLength; pos++ {
        tokenType := lj.Peek(pos).Type
        node, ok = node.children[tokenType]

        if ok && node.leaf {
            largestMatchHandler = node.handler
        }
    }

    if largestMatchHandler == nil {
        // TODO error
        panic(fmt.Sprintf("unexpected token %v", lj.Peek(0).Type))
    }

    // TODO should we assert that the position has changed to prevent infinite loops
    largestMatchHandler(lj, syntax)
}
