package eolparser

import (
	"math"
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/messenger"
	"minimal/minimal-core/built-in/parsers/prattparser"
)

type EOLParser struct {
    prattParser  *prattparser.PrattParser
	eol          lexer.TokenType
	nestingCount uint32
	bindingPower uint32
}

func New(
    messenger *messenger.Messenger,
    prattParser *prattparser.PrattParser,
    lexerScheme *lexer.LexerScheme,
    eol lexer.TokenType,
) *EOLParser {
    lexerScheme.RequireLookahead(2)
    eolParser := EOLParser{prattParser, eol, 0, math.MaxUint32}

    if ok := prattParser.AddPrefix(&eolParser); !ok {
        logEOLPrefixAlreadyDeclared(messenger)
    }

    if ok := prattParser.AddInfix(&eolParser); !ok {
        logEOLInfixAlreadyDeclared(messenger)
    }

    return &eolParser
}

func (e *EOLParser) Parse(l *lexer.Lexer) []ast.Node {
    result := e.prattParser.Parse(l, 0)
    e.bindingPower = math.MaxUint32

    return result
}

func (e *EOLParser) EnterBlock() {
    e.nestingCount++
}

func (e *EOLParser) ExitBlock() {
    if e.nestingCount > 0 {
        e.nestingCount--

        return
    }

    panic("ExitBlock called without matching EnterBlock")
}

func (e *EOLParser) GetTokenType() lexer.TokenType {
    return e.eol
}

func (e *EOLParser) GetBindingPower() uint32 {
    return e.bindingPower
}

func (e *EOLParser) ParsePrefix(pp *prattparser.PrattParser, l *lexer.Lexer, bp uint32) []ast.Node {
    l.Advance()

    return pp.Parse(l, bp)
}

func (e *EOLParser) ParseInfix(pp *prattparser.PrattParser, l *lexer.Lexer, left []ast.Node, bp uint32) []ast.Node {
    if e.nestingCount > 0 {
        l.Advance()

        return left
    }

    next := l.Peek(1).Type

    if pp.IsInfix(next) && !pp.IsPrefix(next) {
        l.Advance()
    } else {
        e.bindingPower = 0
    }

    return left
}

func logEOLPrefixAlreadyDeclared(m *messenger.Messenger) {
    m.Send(messenger.Message{
        Message: "The EOL token has already been declared as a prefix in the pratt parser",
        Severity: messenger.Warning,
        Notes: []string{"The existing prefix parser will be overwritten"},
    })
}

func logEOLInfixAlreadyDeclared(m *messenger.Messenger) {
    m.Send(messenger.Message{
        Message: "The EOL token has already been declared as an infix in the pratt parser",
        Severity: messenger.Warning,
        Notes: []string{"The existing infix parser will be overwritten"},
    })
}
