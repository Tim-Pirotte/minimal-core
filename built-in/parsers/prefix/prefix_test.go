package prefix

import (
	"minimal/minimal-core/built-in/ast"
	"minimal/minimal-core/built-in/lexer"
	symbols "minimal/minimal-core/built-in/matchers/symbol"
	"minimal/minimal-core/built-in/messaging"
	"testing"
)

func TestEmpty(t *testing.T) {
    l := lexer.New(1)
    lj := l.Lex("")

    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    p := NewPrefixParser(m, l, []Rule{})

    p.Parse(lj, ast.New())

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type okParser struct {
    ok bool
}

func (e *okParser) Parse(l *lexer.LexerJob, syntax *ast.AST) {
    e.ok = true
}

func TestNoMatchHandler(t *testing.T) {
    l := lexer.New(1)
    lj := l.Lex("")

    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    m.AddOutput(to)

    ep := okParser{}

    p := NewPrefixParser(m, l, []Rule{{[]lexer.TokenType{}, &ep}})
    p.Parse(lj, ast.New())

    if !ep.ok {
        t.Error("Expected the first handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestAlreadyDeclared(t *testing.T) {
    l := lexer.New(1)
    lj := l.Lex("")

    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
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

    p.Parse(lj, ast.New())

    if firstOk.ok || !secondOk.ok {
        t.Error("Expected only the second handler to run but the first ran")
    }

    if !secondOk.ok {
        t.Error("Expected only the second handler to run but it did not run")
    }

    m.Close()
    to.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "Duplicate prefix in prefix parser",
                Severity: messaging.Error,
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

    l := lexer.New(1)
    a := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
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

    p.Parse(lj, ast.New())

    if !aOk.ok {
        t.Error("Expected the 'a' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !bOk.ok {
        t.Error("Expected the 'b' handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !unknownOk.ok {
        t.Error("Expected the UNKNOWN handler to run")
    }

    lj.Advance()

    p.Parse(lj, ast.New())

    if !emptyOk.ok {
        t.Error("Expected the empty handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestSamePrefix(t *testing.T) {
    aaOk := okParser{}
    abOk := okParser{}

    l := lexer.New(2)
    a := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "a"})
    b := l.NewTokenType(lexer.TokenTypeMetadata{DebugName: "b"})

    sm := symbols.NewSymbolMatcher()
    sm.AddSymbol(l, "a", a)
    sm.AddSymbol(l, "b", b)

    l.AddMatcher(sm)

    m := messaging.NewMessenger()
    to := &messaging.TestOutput{}
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

    p.Parse(lj, ast.New())

    if !aaOk.ok {
        t.Error("Expected the 'aa' handler to run")
    }

    lj.Advance()
    lj.Advance()

    p.Parse(lj, ast.New())

    if !abOk.ok {
        t.Error("Expected the 'ab' handler to run")
    }

    m.Close()
    to.CheckMessages(t, []messaging.Message{})
}
