package lexer

import (
	"bytes"
	"errors"
	"minimal/minimal-lang/built-in/ansi"
	"minimal/minimal-lang/built-in/messenger"
	"minimal/minimal-lang/built-in/outputs/test-output"
	"testing"
)

func TestDisplay(t *testing.T) {
    source := "abc"

    s := NewScheme()
    l := s.Lex(source)

    to := testoutput.TestOutput{}
    m := messenger.New()
    m.AddOutput(&to)

    var buf bytes.Buffer
    d := NewDisplayer(s, &buf, m)

    tokens := []Token{}

    for l.Peek(0).Type != END {
        tokens = append(tokens, l.Peek(0))
        l.Advance()
    }

    d.Display(source, tokens)

    expected := "UNKNOWN              \"a\"                       0..1      (1)\n" +
                "UNKNOWN              \"b\"                       1..2      (1)\n" +
                "UNKNOWN              \"c\"                       2..3      (1)\n"

    actual := buf.String()

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }

    m.Close()
    to.CheckMessages(t, nil)
}

func TestColor(t *testing.T) {
    source := "a"

    s := NewScheme()
    l := s.Lex(source)

    to := testoutput.TestOutput{}
    m := messenger.New()
    m.AddOutput(&to)

    var buf bytes.Buffer
    d := NewDisplayer(s, &buf, m)
    d.SetTokenTypeColor(UNKNOWN, ansi.GetRGBColor(197, 255, 23))

    tokens := []Token{}

    for l.Peek(0).Type != END {
        tokens = append(tokens, l.Peek(0))
        l.Advance()
    }

    d.Display(source, tokens)

    expected := "\x1b[38;2;197;255;23mUNKNOWN             \x1b[0m \"a\"                       0..1      (1)\n"

    actual := buf.String()

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }

    m.Close()
    to.CheckMessages(t, nil)
}

func TestDiff(t *testing.T) {
    source := "ab"

    s := NewScheme()

    l1 := s.Lex(source[:1])
    l2 := s.Lex(source[1:])

    to := testoutput.TestOutput{}
    m := messenger.New()
    m.AddOutput(&to)

    var buf bytes.Buffer
    d := NewDisplayer(s, &buf, m)

    tokens1 := []Token{}

    for l1.Peek(0).Type != END {
        tokens1 = append(tokens1, l1.Peek(0))
        l1.Advance()
    }

    tokens2 := []Token{}

    for l2.Peek(0).Type != END {
        tokens2 = append(tokens2, l2.Peek(0))
        l2.Advance()
    }

    d.DisplayDiff(source, tokens1, tokens2)

    expected := " - UNKNOWN              \"a\"                       0..1      (1)\n" +
                " + UNKNOWN              \"b\"                       1..2      (1)\n"

    actual := buf.String()

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }

    m.Close()
    to.CheckMessages(t, nil)
}

func TestMultiSourceDiff(t *testing.T) {
    source1 := "aa"
    source2 := "ba"

    s := NewScheme()

    l1 := s.Lex(source1)
    l2 := s.Lex(source2)

    to := testoutput.TestOutput{}
    m := messenger.New()
    m.AddOutput(&to)

    var buf bytes.Buffer
    d := NewDisplayer(s, &buf, m)

    tokens1 := []Token{}

    for l1.Peek(0).Type != END {
        tokens1 = append(tokens1, l1.Peek(0))
        l1.Advance()
    }

    tokens2 := []Token{}

    for l2.Peek(0).Type != END {
        tokens2 = append(tokens2, l2.Peek(0))
        l2.Advance()
    }

    d.DisplayMultiSourceDiff(source1, source2, tokens1, tokens2)

    expected := " - UNKNOWN              \"a\"                       0..1      (1)\n" +
                " + UNKNOWN              \"b\"                       0..1      (1)\n" +
                "   UNKNOWN              \"a\"                       1..2      (1)\n"

    actual := buf.String()

    if actual != expected {
        t.Errorf("\nExpected:\n%s\nGot:\n%s", expected, actual)
    }

    m.Close()
    to.CheckMessages(t, nil)
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
    return 0, errors.New("")
}

func TestWriteFail(t *testing.T) {
    source := "a"

    s := NewScheme()
    l := s.Lex(source)

    to := testoutput.New()
    m := messenger.New()
    m.AddOutput(to)

    d := NewDisplayer(s, failingWriter{}, m)

    tokens := []Token{}

    for l.Peek(0).Type != END {
        tokens = append(tokens, l.Peek(0))
        l.Advance()
    }

    d.Display(source, tokens)

    m.Close()
    to.CheckMessages(
        t,
        []messenger.Message{
            {
                Message: "Lexer display output write failed",
                Severity: messenger.Error,
            },
        },
    )
}
