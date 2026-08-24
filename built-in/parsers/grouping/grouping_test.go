package grouping

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	"minimal/minimal-core/built-in/matchers/indentation"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/messenger"
	"minimal/minimal-core/built-in/parsers/pratt"
	"reflect"
	"testing"
)

type binaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint
}

func (b *binaryParser) GetTokenType() lexer.TokenType {
    return b.tokenType
}

func (b *binaryParser) GetBindingPower() uint {
    return b.bindingPower
}

func (b *binaryParser) ParseInfix(
    p *pratt.PrattParser, l *lexer.Lexer, left []ast.Node, minBindingPower uint,
) []ast.Node {
    l.Advance()

    right := p.Parse(l, b.bindingPower)

    result := []ast.Node{{Type: b.nodeType}}
    result = append(result, left...)
    result = append(result, right...)

    return result
}

type unaryParser struct {
    tokenType    lexer.TokenType
    nodeType     ast.NodeType
    bindingPower uint
}

func (u *unaryParser) GetTokenType() lexer.TokenType {
    return u.tokenType
}

func (u *unaryParser) ParsePrefix(
    p *pratt.PrattParser, l *lexer.Lexer, minBindingPower uint,
) []ast.Node {
    l.Advance()

    right := p.Parse(l, u.bindingPower)

    result := []ast.Node{{Type: u.nodeType}}
    result = append(result, right...)

    return result
}

type testGroupingParser struct {
    g *GroupingParser
    l *lexer.LexerScheme
    plus ast.NodeType
    mul ast.NodeType
    min ast.NodeType
    minBin ast.NodeType
    a ast.NodeType
    b ast.NodeType
    c ast.NodeType
}

func getTestGroupingParser() testGroupingParser {
    l := lexer.NewScheme()
    l.RequireLookahead(2)

    openBlockT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "{"})
    closeBlockT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "}"})
    eolT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "{"})
    openParenT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "("})
    closeParenT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: ")"})

    aT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "A"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "B"})
    cT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "C"})

    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "+"})
    mulT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "*"})
    minT := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "-"})

    messenger := messenger.New()
    im := indentation.NewIndentationMatcher(messenger, ':', ' ', openBlockT, closeBlockT, eolT, 0)
    l.AddMatcher(im)

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "b", bT)
    sm.AddSymbol(l, "c", cT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "*", mulT)
    sm.AddSymbol(l, "-", minT)
    sm.AddSymbol(l, "(", openParenT)
    sm.AddSymbol(l, ")", closeParenT)
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
            &unaryParser{minT, min, 2},
        },
        []pratt.Infix{
            &binaryParser{plusT, plus, 2},
            &binaryParser{minT, minBin, 2},
            &binaryParser{mulT, mul, 3},
        },
    )

    g := NewGroupingParser(p, eolT, openParenT, closeParenT)

    return testGroupingParser{g, l, plus, mul, min, minBin, a, b, c}
}

func TestPrefix(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("a*\nb+c")
    result := gp.g.Parse(lj)

    expected := []ast.Node{
        {Type: gp.plus},
        {Type: gp.mul},
        {Type: gp.a},
        {Type: gp.b},
        {Type: gp.c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    lj = gp.l.Lex("a+\nb*c")

    result = gp.g.Parse(lj)

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

func TestInfix(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("a\n*b\n+c")
    result := gp.g.Parse(lj)

    expected := []ast.Node{
        {Type: gp.plus},
        {Type: gp.mul},
        {Type: gp.a},
        {Type: gp.b},
        {Type: gp.c},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }

    lj = gp.l.Lex("a\n+b\n*c")

    result = gp.g.Parse(lj)

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

func TestAmbiguousInfix(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("a\n-b")
    result := gp.g.Parse(lj)
    result = append(result, gp.g.Parse(lj)...)

    expected := []ast.Node{
        {Type: gp.a},
        {Type: gp.min},
        {Type: gp.b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestPrefixWithParentheses(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("(((a)+\n(b)))")
    result := gp.g.Parse(lj)

    expected := []ast.Node{
        {Type: gp.plus},
        {Type: gp.a},
        {Type: gp.b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestInfixWithParentheses(t *testing.T) {
    gp := getTestGroupingParser()

    lj := gp.l.Lex("(a\n-b)")
    result := gp.g.Parse(lj)

    expected := []ast.Node{
        {Type: gp.minBin},
        {Type: gp.a},
        {Type: gp.b},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("\nExpected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestDuplicateEOLPrefix(t *testing.T) {

}

func TestDuplicateEOLInfix(t *testing.T) {

}

func TestDuplicateOpenPrefix(t *testing.T) {

}
