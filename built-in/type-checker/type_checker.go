package typechecker

import "minimal/minimal-core/built-in/ast"

const (
    // Always keep NOTHING at the end
    NOTHING TypeKind = iota
)

type TypeKind uint32

type Type struct {
    Kind      TypeKind
    Reference uint32
}

type TypeKindMetadata interface {
    GetDisplayName(reference uint32) string
}

type TypeCheckerSchema struct {
    astSchema    *ast.ASTSchema
    checkers     map[ast.NodeType]NodeChecker
    lastTypeKind TypeKind
    metadata     []TypeKindMetadata
}

type NodeChecker interface {
    GetNodeType() ast.NodeType
    Check(*TypeChecker) Type
}

type TypeChecker struct {
    astSchema *ast.ASTSchema
    checkers  map[ast.NodeType]NodeChecker
    ast       []ast.Node
    position  uint32
}

func NewSchema(astSchema *ast.ASTSchema) *TypeCheckerSchema {
    return &TypeCheckerSchema{astSchema, map[ast.NodeType]NodeChecker{}, NOTHING, []TypeKindMetadata{}}
}

func (t *TypeCheckerSchema) NewTypeKind(metadata TypeKindMetadata) TypeKind {
    t.lastTypeKind++
    t.metadata = append(t.metadata, metadata)

    return t.lastTypeKind
}

func (t *TypeCheckerSchema) AddNodeChecker(nodeChecker NodeChecker) {
    nodeType := nodeChecker.GetNodeType()

    if _, ok := t.checkers[nodeType]; ok {
        logDuplicateNodeChecker()
    }

    t.checkers[nodeType] = nodeChecker
}

func (t *TypeCheckerSchema) Stringify(dataType Type) string {
    return t.metadata[dataType.Kind].GetDisplayName(dataType.Reference)
}

func (t *TypeCheckerSchema) NewChecker(ast []ast.Node) TypeChecker {
    return TypeChecker{t.astSchema, t.checkers, ast, 0}
}

// TODO what if we want to pass more data around during the semantic pass
func (t *TypeChecker) GetNextType() Type {
    if t.position >= uint32(len(t.ast)) {
        // TODO critical error
    }

    node := t.ast[t.position]
    t.position++

    if checker, ok := t.checkers[node.Type]; ok {
        return checker.Check(t)
    }

    metadata := t.astSchema.GetNodeTypeMetadata(node.Type)

    if metadata.GetChildCount() == ast.VariableChildCount {
        for ; t.ast[t.position].Type != ast.EndNode; t.position++ {
            if t.position + 1 == uint32(len(t.ast)) {
                // TODO error
            }
        }
    } else {
        t.position += uint32(metadata.GetChildCount())

        if t.position >= uint32(len(t.ast)) {
            // TODO error
            // TODO we should stop here since the state is invalid
        }
    }

    return Type{Kind: NOTHING}
}

func (t *TypeChecker) HasReachedEnd() bool {
    return t.position == uint32(len(t.ast))
}

func logDuplicateNodeChecker() {
    // TODO
}
