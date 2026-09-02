package groupingparser

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messenger"
	eolparser "minimal/minimal-core/built-in/parsers/eol"
	"minimal/minimal-core/built-in/parsers/prattparser"
)

type GroupingParser struct {
    m                    *messenger.Messenger
    e                    *eolparser.EOLParser
	open                 lexer.TokenType
	close                lexer.TokenType
    unmatchedOpenMessage string
}

func New(
    m *messenger.Messenger,
    e *eolparser.EOLParser,
    open, close lexer.TokenType,
    unmatchedOpenMessage string,
) *GroupingParser {
    return &GroupingParser{m, e, open, close, unmatchedOpenMessage}
}

func (g *GroupingParser) GetTokenType() lexer.TokenType {
    return g.open
}

func (g *GroupingParser) ParsePrefix(pp *prattparser.PrattParser, l *lexer.Lexer, bp uint32) []ast.Node {
    opening := l.Peek(0)
    l.Advance()

    g.e.EnterBlock()
    result := pp.Parse(l, 0)
    g.e.ExitBlock()

    if l.Peek(0).Type == g.close {
        l.Advance()
    } else {
        g.m.Send(
            messenger.Message{
                Message: g.unmatchedOpenMessage,
                Severity: messenger.Error,
                Context: []messenger.Span{{Content: opening.Value}},
            },
        )
    }

    return result
}
