package prattparser

import (
	"minimal/minimal-lang/built-in/ast"
	"minimal/minimal-lang/built-in/lexer"
	symbols "minimal/minimal-lang/built-in/matchers/symbol"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
    l := lexer.NewScheme()
    p := New(l)
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

func (e *endParser) ParsePrefix(p *PrattParser, l *lexer.Lexer, minBindingPower uint32) []ast.Node {
    return []ast.Node{{Type: e.end}}
}

func TestPrefix(t *testing.T) {
    l := lexer.NewScheme()
    syntax := ast.NewSchema()
    end := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "END"})

    p := New(l)
    p.AddPrefix(&endParser{lexer.END, end})

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

func (p *plusParser) GetBindingPower() uint32 {
    return 2
}

func (p *plusParser) ParseInfix(pp *PrattParser, l *lexer.Lexer, left []ast.Node, bp uint32) []ast.Node {
    l.Advance()

    right := pp.Parse(l, 1)

    result := []ast.Node{{Type: p.plus}}
    result = append(result, left...)
    result = append(result, right...)

    return result
}

func TestBinary(t *testing.T) {
    l := lexer.NewScheme()
    aT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "a", DebugName: "A"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "'+'", DebugName: "+"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "b", DebugName: "B"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "b", bT)
    l.AddMatcher(sm)

    syntax := ast.NewSchema()
    a := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "+", ChildCount: 2})
    b := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "B"})

    p := New(l)
    p.AddPrefix(NewAtomicParser(aT, a))
    p.AddPrefix(NewAtomicParser(bT, b))
    p.AddInfix(&plusParser{plusT, plus})

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

func (m *minusParser) ParsePrefix(p *PrattParser, l *lexer.Lexer, bp uint32) []ast.Node {
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

func (g *groupingParser) ParsePrefix(p *PrattParser, l *lexer.Lexer, bp uint32) []ast.Node {
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

func (e *exclamationParser) GetBindingPower() uint32 {
    return 1
}

func (e *exclamationParser) ParseInfix(
    p *PrattParser, l *lexer.Lexer, left []ast.Node, minBindingPower uint32,
) []ast.Node {
    l.Advance()

    result := []ast.Node{{Type: e.exclamation}}
    result = append(result, left...)

    return result
}

func TestParseCompleteExpression(t *testing.T) {
    l := lexer.NewScheme()
    minusT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "'-'", DebugName: "-"})
    aT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "a", DebugName: "A"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "'+'", DebugName: "+"})
    openParenT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "'('", DebugName: "("})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "b", DebugName: "B"})
    exclamationT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "'!'", DebugName: "!"})
    closeParenT := l.NewTokenType(lexer.TokenTypeMetadata{NounPhrase: "')'", DebugName: ")"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "-", minusT)
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "(", openParenT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "!", exclamationT)
    sm.AddSymbol(l, ")", closeParenT)
    l.AddMatcher(sm)

    syntax := ast.NewSchema()
    minus := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "-", ChildCount: 1})
    a := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "+", ChildCount: 2})
    b := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "B"})
    exclamation := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "!", ChildCount: 1})

    p := New(l)

    p.AddPrefix(&minusParser{minusT, minus})
    p.AddPrefix(NewAtomicParser(aT, a))
    p.AddPrefix(&groupingParser{openParenT, closeParenT, t})
    p.AddPrefix(NewAtomicParser(bT, b))

    p.AddInfix(&plusParser{plusT, plus})
    p.AddInfix(&exclamationParser{exclamationT, exclamation})

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
