package pratt

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
    l := lexer.New(1)
    p := NewPrattParser(l, []Prefix{}, []Infix{})
    lj := l.Lex("")

    p.Parse(lj, 0)
}

type endParser struct {
    endT lexer.TokenType
    end ast.NodeType
}

func (e *endParser) GetTokenType() lexer.TokenType {
    return e.endT
}

func (e *endParser) ParsePrefix(p *PrattParser, l *lexer.LexerJob, minBindingPower uint) []ast.Node {
    return []ast.Node{{Type: e.end}}
}

func TestPrefix(t *testing.T) {
    l := lexer.New(1)
    syntax := ast.New()
    end := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "END"})

    p := NewPrattParser(l,
        []Prefix{&endParser{lexer.END, end}},
        []Infix{},
    )

    lj := l.Lex("")

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: end},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}

type plusParser struct {
    plusT lexer.TokenType
    plus  ast.NodeType
}

func (p *plusParser) GetTokenType() lexer.TokenType {
    return p.plusT
}

func (p *plusParser) GetBindingPower() uint {
    return 2
}

func (p *plusParser) ParseInfix(pp *PrattParser, l *lexer.LexerJob, left []ast.Node, minBindingPower uint) []ast.Node {
    l.Advance()

    right := pp.Parse(l, 1)

    result := []ast.Node{{Type: p.plus}}
    result = append(result, left...)
    result = append(result, right...)

    return result
}

func TestBinary(t *testing.T) {
    l := lexer.New(1)
    aT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "a", DebugName: "A"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'+'", DebugName: "+"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "b", DebugName: "B"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "b", bT)
    l.AddMatcher(sm)

    syntax := ast.New()
    a := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(2, ast.NodeTypeMetadata{DebugName: "+"})
    b := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "B"})

    p := NewPrattParser(
        l,
        []Prefix{NewAtomicParser(aT, a), NewAtomicParser(bT, b)},
        []Infix{&plusParser{plusT, plus}},
    )

    lj := l.Lex("a+b")

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: plus},
        {Type: a},
        {Type: b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}

type minusParser struct {
    minusT lexer.TokenType
    minus  ast.NodeType
}

func (m *minusParser) GetTokenType() lexer.TokenType {
    return m.minusT
}

func (m *minusParser) ParsePrefix(p *PrattParser, l *lexer.LexerJob, minBindingPower uint) []ast.Node {
    l.Advance()

    result := []ast.Node{{Type: m.minus}}
    result = append(result, p.Parse(l, 2)...)

    return result
}

type groupingParser struct {
    openParenT  lexer.TokenType
    closeParenT lexer.TokenType
    t           *testing.T
}

func (g *groupingParser) GetTokenType() lexer.TokenType {
    return g.openParenT
}

func (g *groupingParser) ParsePrefix(p *PrattParser, l *lexer.LexerJob, minBindingPower uint) []ast.Node {
    l.Advance()

    result := p.Parse(l, 0)

    if token := l.Peek(0); token.Type != g.closeParenT {
        g.t.Fatal("Expected ')' matching opening ')'")
    }

    l.Advance()

    return result
}

type exclamationParser struct {
    exclamationT  lexer.TokenType
    exclamation   ast.NodeType
}

func (e *exclamationParser) GetTokenType() lexer.TokenType {
    return e.exclamationT
}

func (e *exclamationParser) GetBindingPower() uint {
    return 1
}

func (e *exclamationParser) ParseInfix(
    p *PrattParser, l *lexer.LexerJob, left []ast.Node, minBindingPower uint,
) []ast.Node {
    l.Advance()

    result := []ast.Node{{Type: e.exclamation}}
    result = append(result, left...)

    return result
}

func TestParseCompleteExpression(t *testing.T) {
    l := lexer.New(1)
    minusT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'-'", DebugName: "-"})
    aT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "a", DebugName: "A"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'+'", DebugName: "+"})
    openParenT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'('", DebugName: "("})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "b", DebugName: "B"})
    exclamationT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'!'", DebugName: "!"})
    closeParenT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "')'", DebugName: ")"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "-", minusT)
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "(", openParenT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "!", exclamationT)
    sm.AddSymbol(l, ")", closeParenT)
    l.AddMatcher(sm)

    syntax := ast.New()
    minus := syntax.NewNodeType(1, ast.NodeTypeMetadata{DebugName: "-"})
    a := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(2, ast.NodeTypeMetadata{DebugName: "+"})
    b := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "B"})
    exclamation := syntax.NewNodeType(1, ast.NodeTypeMetadata{DebugName: "!"})

    p := NewPrattParser(
        l,
        []Prefix{
            &minusParser{minusT, minus},
            NewAtomicParser(aT, a),
            &groupingParser{openParenT, closeParenT, t},
            NewAtomicParser(bT, b),
        },
        []Infix{
            &plusParser{plusT, plus},
            &exclamationParser{exclamationT, exclamation},
        },
    )

    lj := l.Lex("-a+(b!)")

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: plus},
        {Type: minus},
        {Type: a},
        {Type: exclamation},
        {Type: b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }

    lj = l.Lex("-a+b!")

    result = p.Parse(lj, 0)

    expected = []ast.Node{
        {Type: exclamation},
        {Type: plus},
        {Type: minus},
        {Type: a},
        {Type: b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}
