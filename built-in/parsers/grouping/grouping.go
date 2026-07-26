package grouping

import (
	"fmt"
	"math"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

// Start at an end of line token
// Skip it
// Assert that there are no more end of line tokens
// Check if the next token is a led of the pratt parser
// If it isn't also a nud then we continue the pratt parser
// by parsing right with correct binding power propagation

// 5
// * 4

// 5 EOL * 4

// 5
// - 4

// 5 EOL - 4

// 5 *
// 4 + 5

// 5 * EOL 4 + 5

// 5 +
// 4 * 5

// 5 + EOL 4 * 5

type GroupingParser struct {
    prattParser  *pratt.PrattParser
    eolTokenType lexer.TokenType
    bindingPower uint
    nestingCount uint
}

func NewGroupingParser(prattParser *pratt.PrattParser, eolTokenType lexer.TokenType) *GroupingParser {
    g := &GroupingParser{prattParser, eolTokenType, math.MaxUint32, 0}

    if _, ok := prattParser.Prefixes[eolTokenType]; ok {
        logEOLPrefixAlreadyDeclared()
    }

    prattParser.Prefixes[eolTokenType] = g

    if _, ok := prattParser.Infixes[eolTokenType]; ok {
        logEOLInfixAlreadyDeclared()
    }

    prattParser.Infixes[eolTokenType] = g

    return g
}

func (g *GroupingParser) GetTokenType() lexer.TokenType {
    return g.eolTokenType
}

func (g *GroupingParser) GetBindingPower() uint {
    return g.bindingPower
}

func (g *GroupingParser) ParsePrefix(pp *pratt.PrattParser, lj *lexer.LexerJob, bp uint) []ast.Node {
    lj.Advance()

    return pp.Parse(lj, bp)
}

func (g *GroupingParser) ParseInfix(
    pp *pratt.PrattParser,
    lj *lexer.LexerJob,
    left []ast.Node,
    minBindingPower uint,
) []ast.Node {
    if g.nestingCount > 0 {
        lj.Advance()

        return left
    }

    next := lj.Peek(1).Type

    _, isPrefix := pp.Prefixes[next]
    _, isInfix := pp.Infixes[next]

    if isInfix && !isPrefix {
        lj.Advance()
    } else {
        g.bindingPower = 0
    }

    return left
}

func (g *GroupingParser) Parse(l *lexer.LexerJob) []ast.Node {
    result := g.prattParser.Parse(l, 0)

    g.bindingPower = math.MaxUint32

    return result
}

func logEOLPrefixAlreadyDeclared() {
    // TODO use proper message
    fmt.Println("The EOL token has already been used for a prefix in this Pratt parser")
}

func logEOLInfixAlreadyDeclared() {
    // TODO use proper message
    fmt.Println("The EOL token has already been used for an infix in this Pratt parser")
}
