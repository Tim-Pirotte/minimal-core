package typechecker

import "minimal/minimal-lang/built-in/ast"

// Always keep NOTHING at the end
const NOTHING TypeKind = iota

type TypeCheckerSchema struct {
    metadata     []TypeKindMetadata
    lastTypeKind TypeKind
}

type TypeKindMetadata interface {
    GetDisplayName(reference uint32) string
}

type TypeChecker struct {
    typeStack    []TypeFrame
}

type TypeFrame struct {
    incoming Type
    result   Type
}

type Type struct {
    Kind      TypeKind
    Reference uint32
}

type TypeKind uint32

func NewSchema(astSchema *ast.ASTSchema) *TypeCheckerSchema {
    return &TypeCheckerSchema{[]TypeKindMetadata{}, NOTHING}
}

func (t *TypeCheckerSchema) NewTypeKind(metadata TypeKindMetadata) TypeKind {
    t.lastTypeKind++
    t.metadata = append(t.metadata, metadata)

    return t.lastTypeKind
}

func (t *TypeCheckerSchema) Stringify(dataType Type) string {
    return t.metadata[dataType.Kind].GetDisplayName(dataType.Reference)
}

func (t *TypeCheckerSchema) NewChecker(ast []ast.Node) TypeChecker {
    return TypeChecker{[]TypeFrame{{}}}
}

func (t *TypeChecker) GetIncoming() Type {
    if len(t.typeStack) < 2 {
        return Type{Kind: NOTHING}
    }

    return t.typeStack[len(t.typeStack) - 2].incoming
}

func (t *TypeChecker) SetIncoming(incomingType Type) {
    t.typeStack[len(t.typeStack) - 1].incoming = incomingType
}

func (t *TypeChecker) SetResult(resultType Type) {
    if len(t.typeStack) > 1 {
        t.typeStack[len(t.typeStack) - 2].result = resultType
    }
}

func (t *TypeChecker) GetResult() Type {
    return t.typeStack[len(t.typeStack) - 1].result
}

func (t *TypeChecker) Enter() {
    t.typeStack = append(t.typeStack, TypeFrame{Type{Kind: NOTHING}, Type{Kind: NOTHING}})
}

func (t *TypeChecker) Exit() {
    t.typeStack = t.typeStack[:len(t.typeStack) - 1]
}
