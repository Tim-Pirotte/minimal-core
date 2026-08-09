package ast

import (
	"minimal/minimal-core/built-in/messaging"
	"os"
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
    messenger.AddOutput(&messaging.TestOutput{})

    ta.ast.Display(os.Stdout, messenger)
    t.Fail()
}

func TestIncorrectFixedChildren(t *testing.T) {

}

func TestUnclosed(t *testing.T) {

}

func TestNestedUnclosed(t *testing.T) {

}
