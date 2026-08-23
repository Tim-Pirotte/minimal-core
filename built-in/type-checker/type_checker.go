package typechecker

import "minimal/minimal-core/built-in/ast"

const (
    // For generics
    UNRESOLVED TypeKind = iota
    // For things that actually don't return a type like statements
    // Always keep NOTHING at the end
    NOTHING
)

type TypeKind uint32

type Type struct {
    Kind      TypeKind
    Reference uint32
}

type TypeStringify interface {
    Stringify(reference uint32) string
}

type TypeCheckerScheme struct {
    checkers     map[ast.NodeType]NodeChecker
    lastTypeKind TypeKind
    metadata     []TypeStringify
}

type NodeChecker interface {
    Check(*TypeChecker) Type
}

type TypeChecker struct {
    checkers map[ast.NodeType]NodeChecker
    ast      []ast.Node
    position uint32
}

func NewScheme() *TypeCheckerScheme {
    return &TypeCheckerScheme{map[ast.NodeType]NodeChecker{}, NOTHING, []TypeStringify{}}
}

func (t *TypeCheckerScheme) NewTypeKind(stringify TypeStringify) TypeKind {
    t.lastTypeKind++
    t.metadata = append(t.metadata, stringify)

    return t.lastTypeKind
}

func (t *TypeCheckerScheme) AddNodeChecker(nodeType ast.NodeType, nodeChecker NodeChecker) {
    if _, ok := t.checkers[nodeType]; ok {
        logDuplicateNodeChecker()
    }

    t.checkers[nodeType] = nodeChecker
}

func (t *TypeCheckerScheme) Stringify(dataType Type) string {
    return t.metadata[dataType.Kind].Stringify(dataType.Reference)
}

func (t *TypeCheckerScheme) NewChecker(ast []ast.Node) TypeChecker {
    return TypeChecker{t.checkers, ast, 0}
}

func (t *TypeChecker) GetNextType() Type {
    if t.position >= uint32(len(t.ast)) {
        // TODO critical error
    }

    statement := t.ast[t.position]
    t.position++

    if checker, ok := t.checkers[statement.Type]; ok {
        return checker.Check(t)
    }

    return Type{Kind: NOTHING}
}

func (t *TypeChecker) HasReachedEnd() bool {
    return t.position == uint32(len(t.ast))
}

func logDuplicateNodeChecker() {
    // TODO
}
