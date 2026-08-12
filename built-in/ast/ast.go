package ast

const (
    EndNode = NodeType(0)
    VariableChildCount = 255
)

type NodeType uint32

type Node struct {
    Type      NodeType
    Reference uint32
}

type NodeTypeMetadata struct {
    DisplayName string
    DebugName   string
    // Up to 254 fixed children or VariableChildren (255) for a node that ends with EndNode
    ChildCount  uint8
}

type ASTSchema struct {
    lastNodeType NodeType
    metadata     []NodeTypeMetadata
}

type AST []Node

type Traverser struct {
    ast        AST
    position   uint32
}

func New() AST {
    return AST([]Node{})
}

func NewSchema() *ASTSchema {
    return &ASTSchema{EndNode, []NodeTypeMetadata{{DebugName: "EndNode"}}}
}

func (a *ASTSchema) NewNodeType(metadata NodeTypeMetadata) NodeType {
    a.lastNodeType++
    a.metadata = append(a.metadata, metadata)

    return a.lastNodeType
}

func (a *ASTSchema) GetNodeTypeMetadata(nodeType NodeType) NodeTypeMetadata {
    return a.metadata[nodeType]
}

func NewTraverser(ast AST) *Traverser {
    return &Traverser{ast, 0}
}

// Returns nodes depth first pre-order
func (t *Traverser) Next() Node {
    if t.IsAtEnd() {
        panic("*Traveler.Next called when all nodes have been consumed")
    }

    node := t.ast[t.position]
    t.position++

    return node
}

func (t *Traverser) IsAtEnd() bool {
    return int(t.position) == len(t.ast)
}
