package ast

import "math"

type NodeType uint32

type Node struct {
    Type      NodeType
    Reference uint32
}

type NodeTypeMetadata struct {
    DebugName string
}

type AST struct {
    Nodes        []Node
    lastNodeType NodeType
    metadata     []NodeTypeMetadata
}

func New() *AST {
    return &AST{
        []Node{},
        math.MaxUint32,
        []NodeTypeMetadata{},
    }
}

// TODO add child count
func (a *AST) NewNodeType(metadata NodeTypeMetadata) NodeType {
    a.lastNodeType++
    a.metadata = append(a.metadata, metadata)

    return a.lastNodeType
}

func (a *AST) GetNodeTypeMetadata(nodeType NodeType) NodeTypeMetadata {
    return a.metadata[nodeType]
}
