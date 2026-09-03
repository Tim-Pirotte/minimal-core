package prefixunary

import (
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	"minimal/minimal-lang/built-in/parsers/prattparser"
)

type PrefixUnaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint32
}

func NewPrefixUnaryParser(
    tokenType lexer.TokenType, nodeType  ast.NodeType, bindingPower uint32,
) *PrefixUnaryParser {
    return &PrefixUnaryParser{tokenType, nodeType, bindingPower}
}

func (p *PrefixUnaryParser) GetTokenType() lexer.TokenType {
    return p.tokenType
}

func (p *PrefixUnaryParser) ParsePrefix(
    pp *prattparser.PrattParser, l *lexer.Lexer, minBindingPower uint32,
) []ast.Node {
    l.Advance()

    result := []ast.Node{{Type: p.nodeType}}
    result = append(result, pp.Parse(l, p.bindingPower)...)

    return result
}
