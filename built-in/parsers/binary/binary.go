package binary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

type BinaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint32
}

func NewBinaryParser(tokenType lexer.TokenType, nodeType ast.NodeType, bp uint32) *BinaryParser {
    return &BinaryParser{tokenType, nodeType, bp}
}

func (b *BinaryParser) GetTokenType() lexer.TokenType {
    return b.tokenType
}

func (b *BinaryParser) GetBindingPower() uint32 {
    return b.bindingPower
}

func (b *BinaryParser) ParseInfix(p *pratt.PrattParser, l *lexer.Lexer, left []ast.Node, bp uint32) []ast.Node {
    l.Advance()

    right := p.Parse(l, b.bindingPower)

    result := []ast.Node{{Type: b.nodeType}}
    result = append(result, left...)
    result = append(result, right...)

    return result
}
