package ast

import (
    "fmt"
    "io"
    "minimal/minimal-core/built-in/messaging"
    "strings"
)

const spacesPerLevel = 2

type ASTDisplayer struct {
    messenger      *messaging.Messenger
    nodeDisplayers map[NodeType]NodeDisplayer
}

type NodeDisplayer interface {
    GetNodeType() NodeType
    Display(reference uint32) string
}

func NewASTDebugger(messenger *messaging.Messenger) ASTDisplayer {
    return ASTDisplayer{messenger, map[NodeType]NodeDisplayer{}}
}

func (a *ASTDisplayer) AddNodeDisplayer(n NodeDisplayer) {
    nodeType := n.GetNodeType()

    if _, ok := a.nodeDisplayers[nodeType]; ok {
        a.logDuplicateNodeDisplayer(nodeType)
    }

    a.nodeDisplayers[nodeType] = n
}

func (a *ASTDisplayer) Display(ast *AST, o io.Writer) {
    t := NewTraverser(ast)

    for !t.IsAtEnd() {
        node := t.Next()

        if node.Type != EndNode {
            if !a.displayNode(o, t, node, 0) {
                return
            }
        } else if !a.tryWrite(o, "EndNode %s not inside a Node\n", a.getEndNodeName(ast, node)) {
            return
        }
    }
}

// TODO Add syntax highlighting
// Make a color wrapper for displayers and NodeTypes to color them
func (a *ASTDisplayer) displayNode(o io.Writer, t *Traverser, node Node, depth uint) (writeSuccess bool) {
    if !a.tryWrite(
        o,
        "%s%s\n",
        strings.Repeat(" ", int(spacesPerLevel * depth)),
        a.getNodeAsString(t.ast, node),
    ) {
        return false
    }

    childCount := t.ast.GetChildCount(node.Type)

    if childCount == VariableChildren {
        for !t.IsAtEnd() {
            peekedNode := t.ast.Nodes[t.position]

            if peekedNode.Type != EndNode || peekedNode.Reference == uint32(node.Type) {
                node := t.Next()

                if node.Type == EndNode {
                    return true
                }

                if !a.displayNode(o, t, node, depth + 1) {
                    return false
                }
            } else {
                return a.tryWrite(
                    o,
                    "%sIncorrect EndNode %s\n",
                    strings.Repeat(" ", int(spacesPerLevel * depth)),
                    a.getEndNodeName(t.ast, peekedNode),
                )
            }
        }

        return a.tryWrite(o, "%sMissing EndNode\n", strings.Repeat(" ", int(spacesPerLevel * depth)))
    }

    for i := range childCount {
        if t.IsAtEnd() {
            return a.tryWrite(
                o,
                "%s%d missing\n",
                strings.Repeat(" ", int(spacesPerLevel * (depth + 1))),
                childCount - i,
            )
        }

        node := t.Next()

        if node.Type == EndNode {
            if !a.tryWrite(
                o,
                "%sEndNode %s in fixed childcount Node\n",
                strings.Repeat(" ", int(spacesPerLevel * (depth + 1))),
                a.getEndNodeName(t.ast, node),
            ) {
                return false
            }
        } else if !a.displayNode(o, t, node, depth + 1) {
            return false
        }
    }

    return true
}

func (a *ASTDisplayer) tryWrite(output io.Writer, format string, args ...any) bool {
    _, err := fmt.Fprintf(output, format, args...)

    if err != nil {
        a.messenger.Send(
            messaging.Message{
                Message: "AST debugger output write failed",
                Severity: messaging.Error,
            },
        )

        return false
    }

    return true
}

func (a *ASTDisplayer) getEndNodeName(ast *AST, endNode Node) string {
    if int(endNode.Reference) < len(ast.metadata) {
        return ast.GetNodeTypeMetadata(NodeType(endNode.Reference)).DebugName
    }

    return fmt.Sprintf("UNKNOWN Reference=%d", endNode.Reference)
}

func (a *ASTDisplayer) getNodeAsString(ast *AST, node Node) string {
    if displayer, ok := a.nodeDisplayers[node.Type]; ok {
        return displayer.Display(node.Reference)
    }

    metadata := ast.GetNodeTypeMetadata(node.Type)

    if node.Reference == 0 {
        return metadata.DebugName
    }

    return fmt.Sprintf("%s Reference=%d", metadata.DebugName, node.Reference)
}

func (a *ASTDisplayer) logDuplicateNodeDisplayer(nodeType NodeType) {
    // TODO print the debug name
    a.messenger.Send(messaging.Message{
        Message: "Duplicate node displayer in the AST displayer",
        Severity: messaging.Error,
        Notes: []string{fmt.Sprintf("NodeType=%d", nodeType)},
    })
}
