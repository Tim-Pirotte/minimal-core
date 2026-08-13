package prefixunary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

type PrefixUnaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint
}

func NewPrefixUnaryParser(
    tokenType lexer.TokenType, nodeType  ast.NodeType, bindingPower uint,
) *PrefixUnaryParser {
    return &PrefixUnaryParser{tokenType, nodeType, bindingPower}
}

func (p *PrefixUnaryParser) GetTokenType() lexer.TokenType {
    return p.tokenType
}

func (p *PrefixUnaryParser) ParsePrefix(
    pp *pratt.PrattParser, l *lexer.Lexer, minBindingPower uint,
) []ast.Node {
    l.Advance()

    result := []ast.Node{{Type: p.nodeType}}
    result = append(result, pp.Parse(l, p.bindingPower)...)

    return result
}
