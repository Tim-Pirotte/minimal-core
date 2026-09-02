package prattparser

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
)

type Prefix interface {
    GetTokenType() lexer.TokenType
    ParsePrefix(p *PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node
}

type Infix interface {
    GetTokenType() lexer.TokenType
    GetBindingPower() uint32
    ParseInfix(p *PrattParser, lj *lexer.Lexer, left []ast.Node, minBindingPower uint32) []ast.Node
}

type PrattParser struct {
    Prefixes   map[lexer.TokenType]Prefix
    Infixes    map[lexer.TokenType]Infix
}

func New(l *lexer.LexerScheme) *PrattParser {
    return &PrattParser{map[lexer.TokenType]Prefix{}, map[lexer.TokenType]Infix{}}
}

func (p *PrattParser) AddPrefix(prefix Prefix) bool {
    tokenType := prefix.GetTokenType()

    if _, ok := p.Prefixes[tokenType]; ok {
        return false
    }

    p.Prefixes[tokenType] = prefix

    return true
}

func (p *PrattParser) AddInfix(infix Infix) bool {
    tokenType := infix.GetTokenType()

    if _, ok := p.Infixes[tokenType]; ok {
        return false
    }

    p.Infixes[tokenType] = infix

    return true
}

func (p *PrattParser) Parse(l *lexer.Lexer, minBindingPower uint32) []ast.Node {
    prefix, ok := p.Prefixes[l.Peek(0).Type]

    if !ok {
        // TODO proper error message
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

type AtomicParser struct {
	tokenType lexer.TokenType
    nodeType  ast.NodeType
}

func NewAtomicParser(tokenType lexer.TokenType, nodeType ast.NodeType) *AtomicParser {
	return &AtomicParser{tokenType, nodeType}
}

func (a *AtomicParser) GetTokenType() lexer.TokenType {
    return a.tokenType
}

func (a *AtomicParser) ParsePrefix(p *PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
    l.Advance()

    return []ast.Node{{Type: a.nodeType}}
}
