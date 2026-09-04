package ast

import (
	"fmt"
	"io"
	"minimal/minimal-lang/built-in/messenger"
	"strings"
)

const spacesPerLevel = 2

type Displayer struct {
    messenger      *messenger.Messenger
    schema         *ASTSchema
}

func NewDisplayer(messenger *messenger.Messenger, schema *ASTSchema) Displayer {
    return Displayer{messenger, schema}
}

func (a *Displayer) Display(ast []Node, o io.Writer) {
    position := 0

    for position != len(ast) {
        node := ast[position]
        position++

        if node.Type != EndNode {
            if !a.displayNode(o, ast, node, &position, 0) {
                return
            }
        } else if !a.tryWrite(o, "EndNode %s not inside a Node\n", a.getEndNodeName(node)) {
            return
        }
    }
}

func (a *Displayer) DisplayDiff() {

}

func (a *Displayer) displayNode(o io.Writer, ast []Node, node Node, position *int, depth int) (writeSuccess bool) {
    metadata := a.schema.GetNodeTypeMetadata(node.Type)

    if !a.tryWrite(
        o,
        "%s%s\n",
        strings.Repeat(" ", spacesPerLevel * depth),
        metadata.GetDebugName(node.Reference),
    ) {
        return false
    }

    childCount := a.schema.GetNodeTypeMetadata(node.Type).GetChildCount()

    if childCount == VariableChildCount {
        for *position != len(ast) {
            nextNode := ast[*position]

            if nextNode.Type != EndNode || nextNode.Reference == uint32(node.Type) {
                *position++

                if nextNode.Type == EndNode {
                    return true
                }

                if !a.displayNode(o, ast, nextNode, position, depth + 1) {
                    return false
                }
            } else {
                return a.tryWrite(
                    o,
                    "%sIncorrect EndNode %s\n",
                    strings.Repeat(" ", spacesPerLevel * depth),
                    a.getEndNodeName(nextNode),
                )
            }
        }

        return a.tryWrite(o, "%sMissing EndNode\n", strings.Repeat(" ", spacesPerLevel * depth))
    }

    for i := range childCount {
        if *position == len(ast) {
            return a.tryWrite(
                o,
                "%s%d missing\n",
                strings.Repeat(" ", spacesPerLevel * (depth + 1)),
                childCount - i,
            )
        }

        nextNode := ast[*position]
        *position++

        if nextNode.Type == EndNode {
            if !a.tryWrite(
                o,
                "%sEndNode %s in fixed childcount Node\n",
                strings.Repeat(" ", spacesPerLevel * (depth + 1)),
                a.getEndNodeName(nextNode),
            ) {
                return false
            }
        } else if !a.displayNode(o, ast, nextNode, position, depth + 1) {
            return false
        }
    }

    return true
}

func (a *Displayer) tryWrite(output io.Writer, format string, args ...any) bool {
    _, err := fmt.Fprintf(output, format, args...)

    if err != nil {
        a.messenger.Send(
            messenger.Message{
                Message: "AST debugger output write failed",
                Severity: messenger.Error,
            },
        )

        return false
    }

    return true
}

func (a *Displayer) getEndNodeName(endNode Node) string {
    if int(endNode.Reference) < len(a.schema.metadata) {
        return a.schema.GetNodeTypeMetadata(NodeType(endNode.Reference)).GetDebugName(0)
    }

    return fmt.Sprintf("UNKNOWN Reference=%d", endNode.Reference)
}
