package eolparser

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/messenger"
	testoutput "minimal/minimal-core/built-in/outputs/test"
	"minimal/minimal-core/built-in/parsers/binary"
	"minimal/minimal-core/built-in/parsers/prattparser"
	prefixunary "minimal/minimal-core/built-in/parsers/prefix-unary"
	"reflect"
	"testing"
)

type testEOLParser struct {
    e      *EOLParser
    l      *lexer.LexerScheme
    m      *messenger.Messenger
    to     *testoutput.TestOutput
    plus   ast.NodeType
    mul    ast.NodeType
    min    ast.NodeType
    minBin ast.NodeType
    a      ast.NodeType
    b      ast.NodeType
    c      ast.NodeType
}

func getTestEOLParser() testEOLParser {
    l := lexer.NewScheme()

    aT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "A"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "B"})
    cT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "C"})

    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})
    mulT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "*"})
    minT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})
    eolT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "eol"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "c", cT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "*", mulT)
    sm.AddSymbol(l, "-", minT)
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

    p := prattparser.New(l)

    p.AddPrefix(prattparser.NewAtomicParser(aT, a))
    p.AddPrefix(prattparser.NewAtomicParser(bT, b))
    p.AddPrefix(prattparser.NewAtomicParser(cT, c))
    p.AddPrefix(prefixunary.NewPrefixUnaryParser(minT, min, 2))

    p.AddInfix(binary.NewBinaryParser(plusT, plus, 2))
    p.AddInfix(binary.NewBinaryParser(minT, minBin, 2))
    p.AddInfix(binary.NewBinaryParser(mulT, mul, 3))

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    e := New(m, p, l, eolT)

    return testEOLParser{e, l, m, to, plus, mul, min, minBin, a, b, c}
}

func TestPrefix(t *testing.T) {
	eol := getTestEOLParser()

	lj := eol.l.Lex("a*\nb+c")
	result := eol.e.Parse(lj)

	expected := []ast.Node{
		{Type: eol.plus},
		{Type: eol.mul},
		{Type: eol.a},
		{Type: eol.b},
		{Type: eol.c},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\nExpected:\n%+v\nActual:\n%+v", expected, result)
	}

	lj = eol.l.Lex("a+\nb*c")

	result = eol.e.Parse(lj)

	expected = []ast.Node{
		{Type: eol.plus},
		{Type: eol.a},
		{Type: eol.mul},
		{Type: eol.b},
		{Type: eol.c},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
	}

    eol.m.Close()
    eol.to.CheckMessages(t, []messenger.Message{})
}

func TestInfix(t *testing.T) {
	eol := getTestEOLParser()

	lj := eol.l.Lex("a\n*b\n+c")
	result := eol.e.Parse(lj)

	expected := []ast.Node{
		{Type: eol.plus},
		{Type: eol.mul},
		{Type: eol.a},
		{Type: eol.b},
		{Type: eol.c},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
	}

	lj = eol.l.Lex("a\n+b\n*c")

	result = eol.e.Parse(lj)

	expected = []ast.Node{
		{Type: eol.plus},
		{Type: eol.a},
		{Type: eol.mul},
		{Type: eol.b},
		{Type: eol.c},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
	}

    eol.m.Close()
    eol.to.CheckMessages(t, []messenger.Message{})
}

func TestAmbiguousInfix(t *testing.T) {
    eol := getTestEOLParser()

    lj := eol.l.Lex("a\n-b")
    result := eol.e.Parse(lj)
    result = append(result, eol.e.Parse(lj)...)

    expected := []ast.Node{
        {Type: eol.a},
        {Type: eol.min},
        {Type: eol.b},
    }

    // TODO use proper display
    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    eol.m.Close()
    eol.to.CheckMessages(t, []messenger.Message{})
}

func TestAmbiguousInfixInBlock(t *testing.T) {
    eol := getTestEOLParser()

    lj := eol.l.Lex("a\n-b")

    eol.e.EnterBlock()
    result := eol.e.Parse(lj)
    eol.e.ExitBlock()

    expected := []ast.Node{
        {Type: eol.minBin},
        {Type: eol.a},
        {Type: eol.b},
    }

    // TODO use proper display
    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    eol.m.Close()
    eol.to.CheckMessages(t, []messenger.Message{})
}

func TestTooManyExits(t *testing.T) {
    eol := getTestEOLParser()

    defer func() {
        if r := recover(); r != nil {
            if r != "ExitBlock called without matching EnterBlock" {
                t.Errorf(
                    "Expected the panic to be " +
                    "'ExitBlock called without matching EnterBlock' but got '%+v'",
                    r,
                )
            }
        } else {
            t.Errorf("Expected ExitBlock to panic")
        }
    }()

    for range 3 {
        eol.e.EnterBlock()
    }

    for range 4 {
        eol.e.ExitBlock()
    }

    eol.m.Close()
    eol.to.CheckMessages(t, []messenger.Message{})
}

func TestDuplicateEOLPrefixAndInfix(t *testing.T) {
    l := lexer.NewScheme()
    syntax := ast.NewSchema()

    eolT := l.NewTokenType(lexer.TokenTypeMetadata{})
    eol := syntax.NewNodeType(&ast.StructNodeTypeMetadata{})

    p := prattparser.New(l)

    p.AddPrefix(prattparser.NewAtomicParser(eolT, eol))
    p.AddInfix(binary.NewBinaryParser(eolT, eol, 1))

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    New(m, p, l, eolT)

    m.Close()
    to.CheckMessages(t, []messenger.Message{
        {
            Message: "The EOL token has already been declared as a prefix in the pratt parser",
            Severity: messenger.Warning,
            Notes: []string{"The existing prefix parser will be overwritten"},
        },
        {
            Message: "The EOL token has already been declared as an infix in the pratt parser",
            Severity: messenger.Warning,
            Notes: []string{"The existing infix parser will be overwritten"},
        },
    })
}
