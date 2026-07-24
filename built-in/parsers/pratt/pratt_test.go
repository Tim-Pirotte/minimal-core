package pratt

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"reflect"
	"testing"
)

func TestEmpty(t *testing.T) {
    p := NewPrattParser(map[lexer.TokenType]Prefix{}, map[lexer.TokenType]Infix{})

    l := lexer.NewLexer()
    lj := l.Lex("", 1)

    p.Parse(lj, 0)
}

func TestPrefix(t *testing.T) {
    l := lexer.NewLexer()
    syntax := ast.NewAst()
    end := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "END"})

    p := NewPrattParser(
        map[lexer.TokenType]Prefix{
            lexer.END: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    return []ast.Node{{Type: end, Reference: 1}}
                },
            },
        },
        map[lexer.TokenType]Infix{},
    )

    lj := l.Lex("", 1)

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: end, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestBinary(t *testing.T) {
    l := lexer.NewLexer()
    aT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "a", DebugName: "A"})
    plusT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "'+'", DebugName: "+"})
    bT := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "b", DebugName: "B"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", aT)
    sm.AddSymbol(l, "+", plusT)
    sm.AddSymbol(l, "b", bT)
    l.AddMatcher(sm)

    syntax := ast.NewAst()
    a := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "+"})
    b := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "B"})

    p := NewPrattParser(
        map[lexer.TokenType]Prefix{
            aT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    return []ast.Node{{Type: a, Reference: 1}}
                },
            },
            bT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    return []ast.Node{{Type: b, Reference: 1}}
                },
            },
        },
        map[lexer.TokenType]Infix{
            plusT: {
                BindingPower: 1,
                Handler: func(p *PrattParser, lj *lexer.LexerJob, left []ast.Node) []ast.Node {
                    lj.Advance()

                    right := p.Parse(lj, 1)

                    result := []ast.Node{{Type: plus, Reference: 1}}
                    result = append(result, left...)
                    result = append(result, right...)

                    return result
                },
            },
        },
    )

    lj := l.Lex("a+b", 1)

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: plus, Reference: 1},
        {Type: a, Reference: 1},
        {Type: b, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}

func TestParseCompleteExpression(t *testing.T) {
    l := lexer.NewLexer()
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

    syntax := ast.NewAst()
    minus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "-"})
    a := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "A"})
    plus := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "+"})
    b := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "B"})
    exclamation := syntax.NewNodeType(ast.NodeTypeMetadata{DebugName: "!"})

    p := NewPrattParser(
        map[lexer.TokenType]Prefix{
            minusT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    result := []ast.Node{{Type: minus, Reference: 1}}
                    result = append(result, pp.Parse(lj, 2)...)

                    return result
                },
            },
            aT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    return []ast.Node{{Type: a, Reference: 1}}
                },
            },
            openParenT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    result := pp.Parse(lj, 0)

                    if token := lj.Peek(0); token.Type != closeParenT {
                        t.Fatal("Expected ')' matching opening ')'")
                    }

                    lj.Advance()

                    return result
                },
            },
            bT: {
                func(pp *PrattParser, lj *lexer.LexerJob) []ast.Node {
                    lj.Advance()

                    return []ast.Node{{Type: b, Reference: 1}}
                },
            },
        },
        map[lexer.TokenType]Infix{
            plusT: {
                BindingPower: 2,
                Handler: func(p *PrattParser, lj *lexer.LexerJob, left []ast.Node) []ast.Node {
                    lj.Advance()

                    right := p.Parse(lj, 1)

                    result := []ast.Node{{Type: plus, Reference: 1}}
                    result = append(result, left...)
                    result = append(result, right...)

                    return result
                },
            },
            exclamationT: {
                BindingPower: 1,
                Handler: func(p *PrattParser, lj *lexer.LexerJob, left []ast.Node) []ast.Node {
                    lj.Advance()

                    result := []ast.Node{{Type: exclamation, Reference: 1}}
                    result = append(result, left...)

                    return result
                },
            },
        },
    )

    lj := l.Lex("-a+(b!)", 1)

    result := p.Parse(lj, 0)

    expected := []ast.Node{
        {Type: plus, Reference: 1},
        {Type: minus, Reference: 1},
        {Type: a, Reference: 1},
        {Type: exclamation, Reference: 1},
        {Type: b, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }

    lj = l.Lex("-a+b!", 1)

    result = p.Parse(lj, 0)

    expected = []ast.Node{
        {Type: exclamation, Reference: 1},
        {Type: plus, Reference: 1},
        {Type: minus, Reference: 1},
        {Type: a, Reference: 1},
        {Type: b, Reference: 1},
    }

    if !reflect.DeepEqual(expected, result) {
        t.Errorf("Expected:\n%v\nActual:\n%v", expected, result)
    }
}
