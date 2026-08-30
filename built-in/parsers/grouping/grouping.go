package groupingparser

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messenger"
	eolparser "minimal/minimal-core/built-in/parsers/eol"
	"minimal/minimal-core/built-in/parsers/pratt"
)

type GroupingParser struct {
    m     *messenger.Messenger
    eol   *eolparser.EOLParser
	open  lexer.TokenType
	close lexer.TokenType
}

func New(m *messenger.Messenger, eol *eolparser.EOLParser, open, close lexer.TokenType) *GroupingParser {
    return &GroupingParser{m, eol, open, close}
}

func (g *GroupingParser) GetTokenType() lexer.TokenType {
    return g.open
}

func (g *GroupingParser) ParsePrefix(pp *pratt.PrattParser, l *lexer.Lexer, bp uint) []ast.Node {
    opening := l.Peek(0)
    l.Advance()

    g.eol.EnterBlock()
    result := pp.Parse(l, 0)
    g.eol.ExitBlock()

    if l.Peek(0).Type == g.close {
        l.Advance()
    } else {
        // TODO where do we best get the strings from
        g.m.Send(
            messenger.Message{
                Message: "Opening %s does not have a matching closing %s",
                Severity: messenger.Error,
                Context: []messenger.Span{{Content: opening.Value, Note: "The opening %s"}},
            },
        )
    }

    return result
}
