package ast

import "fmt"

const (
    EndNode            = NodeType(0)
    VariableChildCount = 255
)

type NodeType uint32

type Node struct {
    Type      NodeType
    Reference uint32
}

type NodeTypeMetadata interface {
    GetDisplayName(reference uint32) string
    GetDebugName(reference uint32) string
    // Up to 254 fixed children or VariableChildren (255) for a node that ends with EndNode
    GetChildCount() uint8
}

type ASTSchema struct {
    lastNodeType NodeType
    metadata     []NodeTypeMetadata
}

type Traverser struct {
    ast      []Node
    position uint32
}

func NewSchema() *ASTSchema {
    return &ASTSchema{EndNode, []NodeTypeMetadata{&StructNodeTypeMetadata{DebugName: "EndNode"}}}
}

func (a *ASTSchema) NewNodeType(metadata NodeTypeMetadata) NodeType {
    a.lastNodeType++
    a.metadata = append(a.metadata, metadata)

    return a.lastNodeType
}

func (a *ASTSchema) GetNodeTypeMetadata(nodeType NodeType) NodeTypeMetadata {
    return a.metadata[nodeType]
}

type StructNodeTypeMetadata struct {
    DisplayName string
    DebugName   string
    ChildCount  uint8
}

func (s *StructNodeTypeMetadata) GetDisplayName(reference uint32) string {
    return s.DisplayName
}

func (s *StructNodeTypeMetadata) GetDebugName(reference uint32) string {
    if reference == 0 {
        return s.DebugName
    }

    return fmt.Sprintf("%s Reference=%d", s.DebugName, reference)
}

func (s *StructNodeTypeMetadata) GetChildCount() uint8 {
    return s.ChildCount
}
