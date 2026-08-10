package ast

import (
	"bytes"
	"errors"
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

    ta.ast.Display(&buf, messenger)

    expected := `Zero                                                                   (0)
One                                                                    (0)
  Zero                                                                 (0)
Two                                                                    (0)
  One                                                                  (0)
    Zero                                                               (0)
  Variable1                                                            (0)
    Zero                                                               (0)
    Variable2                                                          (0)
      Zero                                                             (0)
    Zero                                                               (0)
`

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

    ta.ast.Display(&buf, messenger)

    expected := `Two                                                                    (0)
  One                                                                  (0)
    1 missing
  1 missing
`

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

    ta.ast.Display(&buf, messenger)

    expected := `Variable1                                                              (0)
  Zero                                                                 (0)
Missing EndNode
`

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

    ta.ast.Display(&buf, messenger)

    expected := `Variable1                                                              (0)
  Zero                                                                 (0)
  Variable2                                                            (0)
    Zero                                                               (0)
    Zero                                                               (0)
  Incorrect EndNode Variable1
`

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

    ta.ast.Display(&buf, messenger)

    expected := `One                                                                    (0)
  EndNode Variable1 in fixed childcount Node
`

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

    ta.ast.Display(&buf, messenger)

    expected := `Variable1                                                              (0)
Incorrect EndNode UNKNOWN (100)
EndNode UNKNOWN (100) not inside a Node
`

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

    ta.ast.Display(writer, messenger)

    messenger.Close()
    to.CheckMessages(t, []messaging.Message{{
        Message: "AST debugger output write failed",
        Severity: messaging.Error,
    }})
}
