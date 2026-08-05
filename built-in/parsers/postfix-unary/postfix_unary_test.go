package postfixunary

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

    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})
    exclamationT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "!"})
    quoteT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "'"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)

    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "!", exclamationT)
    sm.AddSymbol(l, "'", quoteT)

    l.AddMatcher(sm)

    syntax := ast.New()
    plus := syntax.NewNodeType(2, ast.NodeTypeMetadata{DebugName: "+"})
    exclamation := syntax.NewNodeType(1, ast.NodeTypeMetadata{DebugName: "!"})
    quote := syntax.NewNodeType(1, ast.NodeTypeMetadata{DebugName: "'"})

    a := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "a"})
    b := syntax.NewNodeType(0, ast.NodeTypeMetadata{DebugName: "b"})

    p := pratt.NewPrattParser(
        l,
        []pratt.Prefix{
            pratt.NewAtomicParser(aT, a),
            pratt.NewAtomicParser(bT, b),
        },
        []pratt.Infix{
            binary.NewBinaryParser(plusT, plus, 2),
            NewPostfixUnaryParser(exclamationT, exclamation, 2),
            NewPostfixUnaryParser(quoteT, quote, 1),
        },
    )

    lj := l.Lex("a!+b'")
    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: quote},
        {Type: plus},
        {Type: exclamation},
        {Type: a},
        {Type: b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}
