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
    l *lexer.LexerScheme
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
    l := lexer.NewScheme()

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

    syntax := ast.NewSchema()
    plus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "+", ChildCount: 2})
    minus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "-", ChildCount: 1})
    mul := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "*", ChildCount: 2})

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

    lj := bp.l.Lex("a+b")
    result := bp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: bp.plus},
        {Type: bp.a},
        {Type: bp.b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    lj = bp.l.Lex("b+a")
    result = bp.p.Parse(lj, 0)

    expected = []ast.Node{
        {Type: bp.plus},
        {Type: bp.b},
        {Type: bp.a},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestMultiple(t *testing.T) {
    bp := getTestBinaryParser()

    lj := bp.l.Lex("a+b-c*d+e-f")
    result := bp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: bp.minus},
            {Type: bp.plus},
                {Type: bp.minus},
                    {Type: bp.plus},
                        {Type: bp.a},
                        {Type: bp.b},
                    {Type: bp.mul},
                        {Type: bp.c},
                        {Type: bp.d},
                {Type: bp.e},
            {Type: bp.f},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}
