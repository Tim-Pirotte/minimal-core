package groupingparser

import (
    "minimal/minimal-core/built-in/ast"
    "minimal/minimal-core/built-in/lexer"
    symbols "minimal/minimal-core/built-in/matchers/symbol"
    "minimal/minimal-core/built-in/messenger"
    "minimal/minimal-core/built-in/parsers/binary"
    eolparser "minimal/minimal-core/built-in/parsers/eol"
    "minimal/minimal-core/built-in/parsers/pratt"
    prefixunary "minimal/minimal-core/built-in/parsers/prefix-unary"
    "reflect"
    "testing"
)

type testGroupingParser struct {
    g      *GroupingParser
    p      *pratt.PrattParser
    l      *lexer.LexerScheme
    plus   ast.NodeType
    mul    ast.NodeType
    min    ast.NodeType
    minBin ast.NodeType
    a      ast.NodeType
    b      ast.NodeType
    c      ast.NodeType
}

func getTestGroupingParser() testGroupingParser {
    l := lexer.NewScheme()

    openParenT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "("})
    closeParenT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: ")"})

    aT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "A"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "B"})
    cT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "C"})

    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})
    mulT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "*"})
    minT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})

    eolT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "eol"})

    m := messenger.New()

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "c", cT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "*", mulT)
    sm.AddSymbol(l, "-", minT)
    sm.AddSymbol(l, "(", openParenT)
    sm.AddSymbol(l, ")", closeParenT)
    sm.AddSymbol(l, "\n", eolT)
    l.AddMatcher(sm)

    syntax := ast.NewSchema()
    plus := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "+", ChildCount: 2})
    mul := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "*", ChildCount: 2})
    a := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "A"})
    b := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "B"})
    c := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "C"})
    min := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "- (unary)", ChildCount: 1})
    minBin := syntax.NewNodeType(&ast.StructNodeTypeMetadata{DebugName: "- (binary)", ChildCount: 2})

    p := pratt.NewPrattParser(
        l,
        []pratt.Prefix{
            pratt.NewAtomicParser(aT, a),
            pratt.NewAtomicParser(bT, b),
            pratt.NewAtomicParser(cT, c),
            prefixunary.NewPrefixUnaryParser(minT, min, 2),
        },
        []pratt.Infix{
            binary.NewBinaryParser(plusT, plus, 2),
            binary.NewBinaryParser(minT, minBin, 2),
            binary.NewBinaryParser(mulT, mul, 3),
        },
    )

    eol := eolparser.New(m, p, l, eolT)
    g := New(m, eol, openParenT, closeParenT, "'(' does not have a matching ')'")
    p.Prefixes[openParenT] = g

    return testGroupingParser{g, p, l, plus, mul, min, minBin, a, b, c}
}

func TestMultiline(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("((a+\nb)\n*c)")
    result := gp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: gp.mul},
        {Type: gp.plus},
        {Type: gp.a},
        {Type: gp.b},
        {Type: gp.c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    lj = gp.l.Lex("(a+(\nb*\nc))")
    result = gp.p.Parse(lj, 0)

    expected = []ast.Node{
        {Type: gp.plus},
        {Type: gp.a},
        {Type: gp.mul},
        {Type: gp.b},
        {Type: gp.c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestMissingClose(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("((a+b*c)")
    result := gp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: gp.plus},
        {Type: gp.a},
        {Type: gp.mul},
        {Type: gp.b},
        {Type: gp.c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestExtraClose(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("(a+b))")
    result := gp.p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: gp.plus},
        {Type: gp.a},
        {Type: gp.b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}
