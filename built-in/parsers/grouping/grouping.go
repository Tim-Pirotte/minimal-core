package grouping

import (
	"fmt"
	"math"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/parsers/pratt"
)

// TODO should the eol ignore be controllable from outside (for different enclosings)?
// Yes

type GroupingParser struct {
    prattParser  pratt.PrattParser
    eol          eolParser
    nesting      nestingParser
}

type eolParser struct {
    eol lexer.TokenType
    nestingCount *uint
    bindingPower uint
}

type nestingParser struct {
    open  lexer.TokenType
    close lexer.TokenType
    nestingCount *uint
}

func NewGroupingParser(
    prattParser *pratt.PrattParser, eol, open, close lexer.TokenType,
) *GroupingParser {
    nestingCount := uint(0)
    g := &GroupingParser{
        *prattParser,
        eolParser{eol, &nestingCount, math.MaxUint32},
        nestingParser{open, close, &nestingCount},
    }

    if _, ok := prattParser.Prefixes[eol]; ok {
        logEOLPrefixAlreadyDeclared()
    }

    prattParser.Prefixes[eol] = &g.eol

    if _, ok := prattParser.Infixes[eol]; ok {
        logEOLInfixAlreadyDeclared()
    }

    prattParser.Infixes[eol] = &g.eol

    if _, ok := prattParser.Prefixes[open]; ok {
        logOpenPrefixAlreadyDeclared()
    }

    prattParser.Prefixes[open] = &g.nesting

    return g
}

func (e *eolParser) GetTokenType() lexer.TokenType {
    return e.eol
}

func (e *eolParser) GetBindingPower() uint {
    return e.bindingPower
}

func (e *eolParser) ParsePrefix(pp *pratt.PrattParser, l *lexer.LexerJob, bp uint) []ast.Node {
    l.Advance()

    return pp.Parse(l, bp)
}

func (e *eolParser) ParseInfix(
    pp *pratt.PrattParser,
    l *lexer.LexerJob,
    left []ast.Node,
    minBindingPower uint,
) []ast.Node {
    if *e.nestingCount > 0 {
        l.Advance()

        return left
    }

    next := l.Peek(1).Type

    _, isPrefix := pp.Prefixes[next]
    _, isInfix := pp.Infixes[next]

    if isInfix && !isPrefix {
        l.Advance()
    } else {
        e.bindingPower = 0
    }

    return left
}

func (n *nestingParser) GetTokenType() lexer.TokenType {
    return n.open
}

func (n *nestingParser) ParsePrefix(pp *pratt.PrattParser, l *lexer.LexerJob, bp uint) []ast.Node {
    l.Advance()

    *n.nestingCount++
    result := pp.Parse(l, 0)
    *n.nestingCount--

    if l.Peek(0).Type == n.close {
        l.Advance()
    } else {
        // TODO proper message
        fmt.Println("Expected ')' matching the opening '('")
    }

    return result
}

func (g *GroupingParser) Parse(l *lexer.LexerJob) []ast.Node {
    result := g.prattParser.Parse(l, 0)

    g.eol.bindingPower = math.MaxUint32

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

func logOpenPrefixAlreadyDeclared() {
    // TODO use proper message
    fmt.Println("The Open block token has already been used for a prefix in this Pratt parser")
}
