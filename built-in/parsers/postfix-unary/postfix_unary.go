package postfixunary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

type PostfixUnaryParser struct {
	tokenType    lexer.TokenType
	nodeType     ast.NodeType
	bindingPower uint32
}

func NewPostfixUnaryParser(tokenType lexer.TokenType, nodeType ast.NodeType, bindingPower uint32) *PostfixUnaryParser {
	return &PostfixUnaryParser{tokenType, nodeType, bindingPower}
}

func (p *PostfixUnaryParser) GetTokenType() lexer.TokenType {
	return p.tokenType
}

func (p *PostfixUnaryParser) GetBindingPower() uint32 {
	return p.bindingPower
}

func (p *PostfixUnaryParser) ParseInfix(_ *pratt.PrattParser, l *lexer.Lexer, left []ast.Node, bp uint32) []ast.Node {
	l.Advance()

	result := []ast.Node{{Type: p.nodeType}}
	result = append(result, left...)

	return result
}
