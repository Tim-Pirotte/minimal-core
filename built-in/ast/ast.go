package ast

const (
    EndNode = NodeType(0)
    VariableChildren = 255
)

type NodeType uint32

type Node struct {
    Type      NodeType
    Reference uint32
}

type NodeTypeMetadata struct {
    DebugName string
}

// TODO should metadata be in this since it is the same for every source file
type AST struct {
    Nodes        []Node
    childCounts  []uint8
    lastNodeType NodeType
    metadata     []NodeTypeMetadata
}

type Traverser struct {
    ast        *AST
    position   uint32
}

func New() *AST {
    return &AST{
        []Node{},
        []uint8{0},
        EndNode,
        []NodeTypeMetadata{{"EndNode"}},
    }
}

// Create a new node with a certain amount of children (< 255).
// If it has a variable amount of children childCount should be VariableChildren
// and EndNode with reference equal to the opening NodeType should be appended to the AST
// after the last element of the node.
func (a *AST) NewNodeType(childCount uint8, metadata NodeTypeMetadata) NodeType {
    a.lastNodeType++
    a.childCounts = append(a.childCounts, childCount)
    a.metadata = append(a.metadata, metadata)

    return a.lastNodeType
}

func (a *AST) GetChildCount(nodeType NodeType) uint8 {
    return a.childCounts[nodeType]
}

func (a *AST) GetNodeTypeMetadata(nodeType NodeType) NodeTypeMetadata {
    return a.metadata[nodeType]
}

func NewTraverser(ast *AST) *Traverser {
    return &Traverser{ast, 0}
}

// Returns nodes depth first pre-order
// EndNode can be used to know when a variably sized node ends
func (t *Traverser) Next() Node {
    if t.IsAtEnd() {
        panic("*Traveler.Next called when all nodes have been consumed")
    }

    node := t.ast.Nodes[t.position]
    t.position++

    return node
}

func (t *Traverser) IsAtEnd() bool {
    return int(t.position) == len(t.ast.Nodes)
}
