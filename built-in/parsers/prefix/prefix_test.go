package prefix

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"testing"
)

func TestEmpty(t *testing.T) {
    p := NewPrefixParser([]Rule{})

    l := lexer.New()
    lj := l.Lex("", 1)

    p.Parse(lj, ast.New())
}

func TestNoMatchHandler(t *testing.T) {
    ok := false

    p := NewPrefixParser([]Rule{
        {[]lexer.TokenType{}, func(lj *lexer.LexerJob, a *ast.AST) { ok = true }},
    })

    l := lexer.New()
    lj := l.Lex("", 1)

    p.Parse(lj, ast.New())

    if !ok {
        t.Error("Expected the first handler to run")
    }
}

func TestAlreadyDeclared(t *testing.T) {
    ok := false

    p := NewPrefixParser([]Rule{
        {[]lexer.TokenType{}, func(lj *lexer.LexerJob, a *ast.AST) { ok = false }},
        {[]lexer.TokenType{}, func(lj *lexer.LexerJob, a *ast.AST) { ok = true }},
    })

    l := lexer.New()
    lj := l.Lex("", 1)

    p.Parse(lj, ast.New())

    if !ok {
        t.Error("Expected the second handler to run")
    }
}

func TestTokens(t *testing.T) {
    empty_ok := false
    a_ok := false
    b_ok := false
    unknown_ok := false

    l := lexer.New()
    a := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "a", DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "b", DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    p := NewPrefixParser([]Rule{
        {[]lexer.TokenType{}, func(lj *lexer.LexerJob, a *ast.AST) { empty_ok = true }},
        {[]lexer.TokenType{a}, func(lj *lexer.LexerJob, a *ast.AST) { a_ok = true }},
        {[]lexer.TokenType{b}, func(lj *lexer.LexerJob, a *ast.AST) { b_ok = true }},
        {[]lexer.TokenType{lexer.UNKNOWN}, func(lj *lexer.LexerJob, a *ast.AST) { unknown_ok = true }},
    })

    lj := l.Lex("abd", 1)

    p.Parse(lj, ast.New())

    if !a_ok {
        t.Error("Expected the 'a' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !b_ok {
        t.Error("Expected the 'b' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !unknown_ok {
        t.Error("Expected the UNKNOWN handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !empty_ok {
        t.Error("Expected the empty handler to run")
    }
}

func TestSamePrefix(t *testing.T) {
    aa_ok := false
    ab_ok := false

    l := lexer.New()
    a := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "a", DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DisplayName: "b", DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    p := NewPrefixParser([]Rule{
        {[]lexer.TokenType{a, a}, func(lj *lexer.LexerJob, a *ast.AST) { aa_ok = true }},
        {[]lexer.TokenType{a, b}, func(lj *lexer.LexerJob, a *ast.AST) { ab_ok = true }},
    })

    lj := l.Lex("aaab", 2)

    p.Parse(lj, ast.New())

    if !aa_ok {
        t.Error("Expected the 'aa' handler to run")
    }

    lj.Advance()
    lj.Advance()

    p.Parse(lj, ast.New())

    if !ab_ok {
        t.Error("Expected the 'ab' handler to run")
    }
}
