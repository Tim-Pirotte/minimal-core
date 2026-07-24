package pratt

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
)

// TODO every led that is not a nud can be on a new line. This should be disabled when in ()
type Prefix struct {
    Handler      func(*PrattParser, *lexer.LexerJob) []ast.Node
}

type Infix struct {
    BindingPower uint
    Handler      func(p *PrattParser, lj *lexer.LexerJob, left []ast.Node) []ast.Node
}

type PrattParser struct {
    prefixes   map[lexer.TokenType]Prefix
    infixes    map[lexer.TokenType]Infix
}

func NewPrattParser(prefixes map[lexer.TokenType]Prefix, infixes map[lexer.TokenType]Infix) *PrattParser {
    return &PrattParser{prefixes, infixes}
}

func (p *PrattParser) Parse(lj *lexer.LexerJob, minBindingPower uint) []ast.Node {
    prefix, ok := p.prefixes[lj.Peek(0).Type]

    if !ok {
        panic("Expected a valid prefix")
    }

    left := prefix.Handler(p, lj)
    infix, ok := p.infixes[lj.Peek(0).Type]

    for ok && infix.BindingPower > minBindingPower {
        left = infix.Handler(p, lj, left)
        infix, ok = p.infixes[lj.Peek(0).Type]
    }

    return left
}
