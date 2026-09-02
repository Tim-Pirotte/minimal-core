package postfixunary

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/parsers/binary"
	"minimal/minimal-core/built-in/parsers/prattparser"
	"reflect"
	"testing"
)

func TestBindingPower(t *testing.T) {
    l := lexer.NewScheme()

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

    syntax := ast.NewSchema()
    plus := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "+", ChildCount: 2})
    exclamation := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "!", ChildCount: 1})
    quote := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "'", ChildCount: 1})

    a := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "a"})
    b := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "b"})

    p := prattparser.New(l)

    p.AddPrefix(prattparser.NewAtomicParser(aT, a))
    p.AddPrefix(prattparser.NewAtomicParser(bT, b))

    p.AddInfix(binary.NewBinaryParser(plusT, plus, 2))
    p.AddInfix(NewPostfixUnaryParser(exclamationT, exclamation, 2))
    p.AddInfix(NewPostfixUnaryParser(quoteT, quote, 1))

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
