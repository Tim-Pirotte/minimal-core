package ast

import (
	"bytes"
	"errors"
	"minimal/minimal-core/built-in/ansi"
	"minimal/minimal-core/built-in/messenger"
	"testing"
)

type testAST struct {
    a                      ASTDisplayer
    schema                 *ASTSchema
    messenger              *messenger.Messenger
    to                     *messenger.TestOutput
    zeroChildren           NodeType
    oneChild               NodeType
    twoChildren            NodeType

    firstVariableChildren  NodeType
    secondVariableChildren NodeType
}

func getTestAST() testAST {
    schema := NewSchema()

    zeroChildren := schema.NewNodeType(NodeTypeMetadata{DebugName: "Zero"})
    oneChild := schema.NewNodeType(NodeTypeMetadata{DebugName: "One", ChildCount: 1})
    twoChildren := schema.NewNodeType(NodeTypeMetadata{DebugName: "Two", ChildCount: 2})

    firstVariableChildren := schema.NewNodeType(
        NodeTypeMetadata{DebugName: "Variable1", ChildCount: VariableChildCount},
    )

    secondVariableChildren := schema.NewNodeType(
        NodeTypeMetadata{DebugName: "Variable2", ChildCount: VariableChildCount},
    )

    m := messenger.New()
    to := &messenger.TestOutput{}
    m.AddOutput(to)

    a := NewASTDisplayer(m, schema)

    return testAST{
        a,
        schema,
        m,
        to,
        zeroChildren,
        oneChild,
        twoChildren,
        firstVariableChildren,
        secondVariableChildren,
    }
}

func TestCorrect(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
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

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

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

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestIncorrectFixedChildren(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.twoChildren, 0},
        {ta.oneChild, 0},
    }

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

    expected := "Two\n" +
                "  One\n" +
                "    1 missing\n" +
                "  1 missing\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestMissingEndNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
    }

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

    expected := "Variable1\n" +
                "  Zero\n" +
                "Missing EndNode\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestMissingEndNodeNested(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.secondVariableChildren, 0},
        {ta.zeroChildren, 0},
        {ta.zeroChildren, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

    expected := "Variable1\n" +
                "  Zero\n" +
                "  Variable2\n" +
                "    Zero\n" +
                "    Zero\n" +
                "  Incorrect EndNode Variable1\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestEndNodeInFixedChildrenNode(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.oneChild, 0},
        {EndNode, uint32(ta.firstVariableChildren)},
    }

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

    expected := "One\n" +
                "  EndNode Variable1 in fixed childcount Node\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestUnknownEndNodeReference(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.firstVariableChildren, 0},
        {EndNode, uint32(100)},
    }

    var buf bytes.Buffer

    ta.a.Display(ast, &buf)

    expected := "Variable1\n" +
                "Incorrect EndNode UNKNOWN Reference=100\n" +
                "EndNode UNKNOWN Reference=100 not inside a Node\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

type failingWriter struct{}

func (failingWriter) Write(p []byte) (int, error) {
    return 0, errors.New("")
}

func TestFailingWriter(t *testing.T) {
    ta := getTestAST()
    writer := failingWriter{}

    ast := []Node{{ta.zeroChildren, 0}}

    ta.a.Display(ast, writer)

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{{
        Message: "AST debugger output write failed",
        Severity: messenger.Error,
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

    ast := []Node{{ta.oneChild, 0}, {ta.zeroChildren, 0}}

    var buf bytes.Buffer

    cd := customDisplay{ta.zeroChildren, "TestCustomDisplay"}
    ta.a.AddNodeDisplayer(&cd)
    ta.a.Display(ast, &buf)

    expected := "One\n" +
                "  TestCustomDisplay\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}

func TestDuplicateDisplay(t *testing.T) {
    ta := getTestAST()
    ast := []Node{{ta.zeroChildren, 0}}

    var buf bytes.Buffer

    cd1 := customDisplay{ta.zeroChildren, "FirstCustomDisplay"}
    cd2 := customDisplay{ta.zeroChildren, "SecondCustomDisplay"}

    ta.a.AddNodeDisplayer(&cd1)
    ta.a.AddNodeDisplayer(&cd2)

    ta.a.Display(ast, &buf)

    expected := "SecondCustomDisplay\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(
        t,
        []messenger.Message{
            {
                Message: "Duplicate node displayer in the AST displayer",
                Severity: messenger.Error,
                Notes: []string{"NodeType=Zero"},
            },
        },
    )
}

func TestColoredNodes(t *testing.T) {
    ta := getTestAST()

    ast := []Node{
        {ta.oneChild, 0},
        {ta.zeroChildren, 1},
    }

    var buf bytes.Buffer

    ta.a.AddNodeDisplayer(NewNodeColorer(ta.schema, ta.oneChild, ansi.GetRGBColor(0, 255, 242)))
    ta.a.AddNodeDisplayer(NewNodeColorer(ta.schema, ta.zeroChildren, ansi.GetRGBColor(255, 0, 162)))
    ta.a.Display(ast, &buf)

    expected := "\x1b[38;2;0;255;242mOne\x1b[0m\n" +
                "  \x1b[38;2;255;0;162mZero\x1b[0m Reference=1\n"

    if buf.String() != expected {
        t.Errorf("\nExpected:\n%sGot:\n%s", expected, buf.String())
    }

    ta.messenger.Close()
    ta.to.CheckMessages(t, []messenger.Message{})
}
