package ast

import (
	"bytes"
	"errors"
	"minimal/minimal-core/built-in/ansi"
	"minimal/minimal-core/built-in/messaging"
	"testing"
)

type testAST struct {
    ast *AST
    zeroChildren NodeType
    oneChild NodeType
    twoChildren NodeType

    firstVariableChildren NodeType
    secondVariableChildren NodeType
}

func getTestAST() testAST {
    ast := New()

    zeroChildren := ast.NewNodeType(0, NodeTypeMetadata{DebugName: "Zero"})
    oneChild := ast.NewNodeType(1, NodeTypeMetadata{DebugName: "One"})
    twoChildren := ast.NewNodeType(2, NodeTypeMetadata{DebugName: "Two"})

    firstVariableChildren := ast.NewNodeType(VariableChildren, NodeTypeMetadata{DebugName: "Variable1"})
    secondVariableChildren := ast.NewNodeType(VariableChildren, NodeTypeMetadata{DebugName: "Variable2"})

    return testAST{ast, zeroChildren, oneChild, twoChildren, firstVariableChildren, secondVariableChildren}
}

func TestCorrect(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.zeroChildren, 0},
        {ta.oneChild, 0},
        {ta.zeroChildren, 0},
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
        {ta.zeroChildren, 0},
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.secondVariableChildren, 0},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.secondVariableChildren)},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "Zero\n" +
                "One\n" +
                "  Zero\n" +
                "Two\n" +
                "  One\n" +
                "    Zero\n" +
                "  Variable1\n" +
                "    Zero\n" +
                "    Variable2\n" +
                "      Zero\n" +
                "    Zero\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestIncorrectFixedChildren(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "Two\n" +
                "  One\n" +
                "    1 missing\n" +
                "  1 missing\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestMissingEndNode(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "Variable1\n" +
                "  Zero\n" +
                "Missing EndNode\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestMissingEndNodeNested(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.secondVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "Variable1\n" +
                "  Zero\n" +
                "  Variable2\n" +
                "    Zero\n" +
                "    Zero\n" +
                "  Incorrect EndNode Variable1\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestEndNodeInFixedChildrenNode(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.oneChild, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "One\n" +
                "  EndNode Variable1 in fixed childcount Node\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestUnknownEndNodeReference(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.firstVariableChildren, 0},
        {EndNode, uint32(100)},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, &buf)

    expected := "Variable1\n" +
                "Incorrect EndNode UNKNOWN Reference=100\n" +
                "EndNode UNKNOWN Reference=100 not inside a Node\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
    return 0, errors.New("")
}

func TestFailingWriter(t *testing.T) {
    ta := getTestAST()
    writer := failingWriter{}
    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    ta.ast.Nodes = []Node{{ta.zeroChildren, 0}}

    a := NewASTDebugger(messenger)
    a.Display(ta.ast, writer)

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{{
        Message: "AST debugger output write failed",
        Severity: messaging.Error,
    }})
}

type customDisplay struct {
    nodeType NodeType
    message  string
}

func (c *customDisplay) GetNodeType() NodeType {
    return c.nodeType
}

func (c *customDisplay) Display(uint32) string {
    return c.message
}

func TestCustomDisplay(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{{ta.oneChild, 0}, {ta.zeroChildren, 0}}

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    cd := customDisplay{ta.zeroChildren, "TestCustomDisplay"}
    a.AddNodeDisplayer(&cd)
    a.Display(ta.ast, &buf)

    expected := "One\n" +
                "  TestCustomDisplay\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}

func TestDuplicateDisplay(t *testing.T) {
    ta := getTestAST()
    ta.ast.Nodes = []Node{{ta.zeroChildren, 0}}

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)

    cd1 := customDisplay{ta.zeroChildren, "FirstCustomDisplay"}
    cd2 := customDisplay{ta.zeroChildren, "SecondCustomDisplay"}

    a.AddNodeDisplayer(&cd1)
    a.AddNodeDisplayer(&cd2)

    a.Display(ta.ast, &buf)

    expected := "SecondCustomDisplay\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(
        t,
        []messaging.Message{
            {
                Message: "Duplicate node displayer in the AST displayer",
                Severity: messaging.Error,
                Notes: []string{"NodeType=1"},
            },
        },
    )
}

func TestColoredNodes(t *testing.T) {
    ta := getTestAST()

    ta.ast.Nodes = []Node{
        {ta.oneChild, 0},
        {ta.zeroChildren, 1},
    }

    messenger := messaging.NewMessenger()
    to := &messaging.TestOutput{}
    messenger.AddOutput(to)

    var buf bytes.Buffer

    a := NewASTDebugger(messenger)
    a.AddNodeDisplayer(NewNodeColorer(ta.ast, ta.oneChild, ansi.GetRGBColor(0, 255, 242)))
    a.AddNodeDisplayer(NewNodeColorer(ta.ast, ta.zeroChildren, ansi.GetRGBColor(255, 0, 162)))
    a.Display(ta.ast, &buf)

    expected := "\x1b[38;2;0;255;242mOne\x1b[0m\n" +
                "  \x1b[38;2;255;0;162mZero\x1b[0m Reference=1\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{})
}
