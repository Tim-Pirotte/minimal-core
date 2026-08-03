package prefixunary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/parsers/binary"
	"minimal/minimal-core/built-in/parsers/pratt"
	"reflect"
	"testing"
)

func TestBindingPower(t *testing.T) {
    l := lexer.New(1)

    aT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})
    cT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "c"})
    dT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "d"})

    minusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})
    minusMinusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "c", cT)
    sm.AddSymbol(l, "d", dT)

    sm.AddSymbol(l, "-", minusT)
    sm.AddSymbol(l, "--", minusMinusT)
    sm.AddSymbol(l, "+", plusT)

    l.AddMatcher(sm)

    syntax := ast.New()
    minus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "-"})
    minusMinus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "--"})
    plus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "+"})

    a := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "a"})
    b := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "b"})
    c := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "c"})
    d := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "d"})

    p := pratt.NewPrattParser(
        l,
        []pratt.Prefix{
            pratt.NewAtomicParser(aT, a),
            pratt.NewAtomicParser(bT, b),
            pratt.NewAtomicParser(cT, c),
            pratt.NewAtomicParser(dT, d),
            NewPrefixUnaryParser(minusT, minus, 2),
            NewPrefixUnaryParser(minusMinusT, minusMinus, 1),
        },
        []pratt.Infix{
            binary.NewBinaryParser(plusT, plus, 2),
        },
    )

    lj := l.Lex("-a+--b+c")
    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: plus},
        {Type: minus},
        {Type: a},
        {Type: minusMinus},
        {Type: plus},
        {Type: b},
        {Type: c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}
