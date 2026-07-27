package binary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

type BinaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint
}

func NewBinaryParser(tokenType lexer.TokenType, nodeType ast.NodeType, bindingPower uint) *BinaryParser {
    return &BinaryParser{tokenType, nodeType, bindingPower}
}

func (b *BinaryParser) GetTokenType() lexer.TokenType {
    return b.tokenType
}

func (b *BinaryParser) GetBindingPower() uint {
    return b.bindingPower
}

func (b *BinaryParser) ParseInfix(
    p *pratt.PrattParser, l *lexer.LexerJob, left []ast.Node, minBindingPower uint,
) []ast.Node {
    l.Advance()

    right := p.Parse(l, b.bindingPower)

    result := []ast.Node{{Type: b.nodeType}}
    result = append(result, left...)
    result = append(result, right...)

    return result
}
