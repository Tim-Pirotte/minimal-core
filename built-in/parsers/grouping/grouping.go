package grouping

import (
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
    eolTokenType lexer.TokenType
    nestingCount uint
}

func NewGroupingParser(eolTokenType lexer.TokenType) *GroupingParser {
    return &GroupingParser{eolTokenType, 0}
}

func (g *GroupingParser) ParsePrefix(pp *pratt.PrattParser, lj *lexer.LexerJob, bp uint) []ast.Node {
    lj.Advance()

    return pp.Parse(lj, bp)
}

func (g *GroupingParser) ParseInfix(
    pp *pratt.PrattParser,
    lj *lexer.LexerJob,
    left []ast.Node,
) []ast.Node {
    // This should have a binding power of 2^32 - 1
    // The first time this runs if we don't advance
    // then the bp will be the same and the loop will stop.
    // This does mean that every Parse call to the pratt parser will visit this function once sadly
    next := lj.Peek(1).Type

    _, isPrefix := pp.Prefixes[next]
    _, isInfix := pp.Infixes[next]

    if isInfix && !isPrefix {
        lj.Advance()
    }

    return left
}
