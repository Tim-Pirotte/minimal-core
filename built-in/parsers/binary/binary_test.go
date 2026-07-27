package binary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/parsers/pratt"
	"reflect"
	"testing"
)

type testBinaryParser struct {
    p *pratt.PrattParser
    l *lexer.Lexer
    a ast.NodeType
    b ast.NodeType
    c ast.NodeType
    d ast.NodeType
    e ast.NodeType
    f ast.NodeType
    plus ast.NodeType
    minus ast.NodeType
    mul ast.NodeType
}

func getTestBinaryParser() testBinaryParser {
    l := lexer.NewLexer()

    aT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})
    cT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "c"})
    dT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "d"})
    eT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "e"})
    fT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "f"})

    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})
    minusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})
    mulT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "*"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "c", cT)
    sm.AddSymbol(l, "d", dT)
    sm.AddSymbol(l, "e", eT)
    sm.AddSymbol(l, "f", fT)

    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "-", minusT)
    sm.AddSymbol(l, "*", mulT)
    l.AddMatcher(sm)

    syntax := ast.NewAst()
    plus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "+"})
    minus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "-"})
    mul := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "*"})

    a := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "a"})
    b := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "b"})
    c := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "c"})
    d := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "d"})
    e := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "e"})
    f := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "f"})

    p := pratt.NewPrattParser(
        l,
        []pratt.Prefix{
            pratt.NewAtomicParser(aT, a),
            pratt.NewAtomicParser(bT, b),
            pratt.NewAtomicParser(cT, c),
            pratt.NewAtomicParser(dT, d),
            pratt.NewAtomicParser(eT, e),
            pratt.NewAtomicParser(fT, f),
        },
        []pratt.Infix{
            NewBinaryParser(plusT, plus, 2),
            NewBinaryParser(minusT, minus, 2),
            NewBinaryParser(mulT, mul, 3),
        },
    )

    return testBinaryParser{p, l, a, b, c, d, e, f, plus, minus, mul}
}

func TestSingle(t *testing.T) {
    bp := getTestBinaryParser()

    lj := bp.l.Lex("a+b", 1)
    result := bp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: bp.plus, Reference: 1},
        {Type: bp.a, Reference: 1},
        {Type: bp.b, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    lj = bp.l.Lex("b+a", 1)
    result = bp.p.Parse(lj, 0)

    expected = []ast.Node{
        {Type: bp.plus, Reference: 1},
        {Type: bp.b, Reference: 1},
        {Type: bp.a, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestMultiple(t *testing.T) {
    bp := getTestBinaryParser()

    lj := bp.l.Lex("a+b-c*d+e-f", 1)
    result := bp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: bp.minus, Reference: 1},
            {Type: bp.plus, Reference: 1},
                {Type: bp.minus, Reference: 1},
                    {Type: bp.plus, Reference: 1},
                        {Type: bp.a, Reference: 1},
                        {Type: bp.b, Reference: 1},
                    {Type: bp.mul, Reference: 1},
                        {Type: bp.c, Reference: 1},
                        {Type: bp.d, Reference: 1},
                {Type: bp.e, Reference: 1},
            {Type: bp.f, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}
