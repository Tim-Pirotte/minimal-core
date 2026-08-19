package prefix

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/messenger"
	testoutput "minimal/minimal-core/built-in/outputs/test"
	"testing"
)

func TestEmpty(t *testing.T) {
    l := lexer.NewScheme()
    lj := l.Lex("")

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    p := NewPrefixParser(m, l, []Rule{})

    p.Parse(lj, ast.NewSchema())

    m.Close()
    to.CheckMessages(t, []messenger.Message{})
}

type okParser struct {
    ok bool
}

func (e *okParser) Parse(l *lexer.Lexer, syntax *ast.ASTSchema) {
    e.ok = true
}

func TestNoMatchHandler(t *testing.T) {
    l := lexer.NewScheme()
    lj := l.Lex("")

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    ep := okParser{}

    p := NewPrefixParser(m, l, []Rule{{[]lexer.TokenType{}, &ep}})
    p.Parse(lj, ast.NewSchema())

    if !ep.ok {
        t.Error("Expected the first handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messenger.Message{})
}

func TestAlreadyDeclared(t *testing.T) {
    l := lexer.NewScheme()
    lj := l.Lex("")

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    firstOk := &okParser{}
    secondOk := &okParser{}

    p := NewPrefixParser(
        m,
        l,
        []Rule{
            {[]lexer.TokenType{}, firstOk},
            {[]lexer.TokenType{}, secondOk},
        },
    )

    p.Parse(lj, ast.NewSchema())

    if firstOk.ok || !secondOk.ok {
        t.Error("Expected only the second handler to run but the first ran")
    }

    if !secondOk.ok {
        t.Error("Expected only the second handler to run but it did not run")
    }

    m.Close()
    to.CheckMessages(
        t,
        []messenger.Message{
            {
                Message: "Duplicate prefix in prefix parser",
                Severity: messenger.Error,
                Notes: []string{"Prefix: []"},
            },
        },
    )
}

func TestTokens(t *testing.T) {
    emptyOk := okParser{}
    aOk := okParser{}
    bOk := okParser{}
    unknownOk := okParser{}

    l := lexer.NewScheme()
    a := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    p := NewPrefixParser(
        m,
        l,
        []Rule{
            {[]lexer.TokenType{}, &emptyOk},
            {[]lexer.TokenType{a}, &aOk},
            {[]lexer.TokenType{b}, &bOk},
            {[]lexer.TokenType{lexer.UNKNOWN}, &unknownOk},
        },
    )

    lj := l.Lex("abd")

    p.Parse(lj, ast.NewSchema())

    if !aOk.ok {
        t.Error("Expected the 'a' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.NewSchema())

    if !bOk.ok {
        t.Error("Expected the 'b' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.NewSchema())

    if !unknownOk.ok {
        t.Error("Expected the UNKNOWN handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.NewSchema())

    if !emptyOk.ok {
        t.Error("Expected the empty handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messenger.Message{})
}

func TestSamePrefix(t *testing.T) {
    aaOk := okParser{}
    abOk := okParser{}

    l := lexer.NewScheme()
    l.RequireLookahead(2)

    a := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    m := messenger.New()
    to := testoutput.New()
    m.AddOutput(to)

    p := NewPrefixParser(
        m,
        l,
        []Rule{
            {[]lexer.TokenType{a, a}, &aaOk},
            {[]lexer.TokenType{a, b}, &abOk},
        },
    )

    lj := l.Lex("aaab")

    p.Parse(lj, ast.NewSchema())

    if !aaOk.ok {
        t.Error("Expected the 'aa' handler to run")
    }

    lj.Advance()
    lj.Advance()

    p.Parse(lj, ast.NewSchema())

    if !abOk.ok {
        t.Error("Expected the 'ab' handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messenger.Message{})
}
