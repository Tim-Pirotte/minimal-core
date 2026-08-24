package ast

import (
	"bytes"
	"errors"
	"minimal/minimal-core/built-in/messenger"
	testoutput "minimal/minimal-core/built-in/outputs/test"
	"testing"
)

type testAST struct {
    a                      ASTDisplayer
    schema                 *ASTSchema
    messenger              *messenger.Messenger
    to                     *testoutput.TestOutput
    zeroChildren           NodeType
    oneChild               NodeType
    twoChildren            NodeType

    firstVariableChildren  NodeType
    secondVariableChildren NodeType
}

func getTestAST() testAST {
    schema := NewSchema()

    zeroChildren := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "Zero"})
    oneChild := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "One", ChildCount: 1})
    twoChildren := schema.NewNodeType(&StructNodeTypeMetadata{DebugName: "Two", ChildCount: 2})

    firstVariableChildren := schema.NewNodeType(
        &StructNodeTypeMetadata{DebugName: "Variable1", ChildCount: VariableChildCount},
    )

    secondVariableChildren := schema.NewNodeType(
        &StructNodeTypeMetadata{DebugName: "Variable2", ChildCount: VariableChildCount},
    )

    m := messenger.New()
    to := testoutput.New()
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
