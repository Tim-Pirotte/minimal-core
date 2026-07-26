package pratt

import (
	"fmt"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
)

type Prefix interface {
    GetTokenType() lexer.TokenType
    ParsePrefix(p *PrattParser, l *lexer.LexerJob, minBindingPower uint) []ast.Node
}

type Infix interface {
    GetTokenType() lexer.TokenType
    GetBindingPower() uint
    ParseInfix(p *PrattParser, lj *lexer.LexerJob, left []ast.Node, minBindingPower uint) []ast.Node
}

type PrattParser struct {
    Prefixes   map[lexer.TokenType]Prefix
    Infixes    map[lexer.TokenType]Infix
}

func NewPrattParser(l *lexer.Lexer, prefixes []Prefix, infixes []Infix) *PrattParser {
    prefixMap := map[lexer.TokenType]Prefix{}

    for _, prefix := range prefixes {
        tokenType := prefix.GetTokenType()

        if _, ok := prefixMap[tokenType]; ok {
            logDuplicatePrefix(l, tokenType)
        } else {
            prefixMap[tokenType] = prefix
        }
    }

    infixMap := map[lexer.TokenType]Infix{}

    for _, infix := range infixes {
        tokenType := infix.GetTokenType()

        if _, ok := infixMap[tokenType]; ok {
            logDuplicateInfix(l, tokenType)
        } else {
            infixMap[tokenType] = infix
        }
    }

    return &PrattParser{prefixMap, infixMap}
}

func (p *PrattParser) Parse(l *lexer.LexerJob, minBindingPower uint) []ast.Node {
    prefix, ok := p.Prefixes[l.Peek(0).Type]

    if !ok {
        panic("Expected a valid prefix")
    }

    left := prefix.ParsePrefix(p, l, minBindingPower)
    infix, ok := p.Infixes[l.Peek(0).Type]

    for ok && infix.GetBindingPower() > minBindingPower {
        left = infix.ParseInfix(p, l, left, minBindingPower)
        infix, ok = p.Infixes[l.Peek(0).Type]
    }

    return left
}

func logDuplicatePrefix(l *lexer.Lexer, tokenType lexer.TokenType) {
    // TODO proper message
    fmt.Printf(
        "Multiple prefixes have been declared for the token type %s\n",
        l.GetTokenTypeMetadata(tokenType),
    )
}

func logDuplicateInfix(l *lexer.Lexer, tokenType lexer.TokenType) {
    // TODO proper message
    fmt.Printf(
        "Multiple infixes have been declared for the token type %s\n",
        l.GetTokenTypeMetadata(tokenType),
    )
}
