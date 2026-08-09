package ast

import (
    "fmt"
    "io"
    "minimal/minimal-core/built-in/messaging"
    "strings"
)

const spacesPerLevel = 2

type astDebugger struct {
    traverser       *Traverser
    output         io.Writer
    messenger      *messaging.Messenger
}

func (a *AST) Display(output io.Writer, messenger *messaging.Messenger) {
    ad := &astDebugger{NewTraverser(a), output, messenger}

    for !ad.traverser.IsAtEnd() {
        node := ad.traverser.Next()
        ad.displayNode(node, 0)
    }
}

func (a *astDebugger) displayNode(node Node, depth uint) (writeSuccess bool) {
    metadata := a.traverser.ast.GetNodeTypeMetadata(node.Type)

    // TODO allow custom writes for references
    if !a.tryWrite(
        "%-70s (%d)\n",
        strings.Repeat(" ", int(spacesPerLevel * depth)) + metadata.DebugName,
        node.Reference,
    ) {
        return false
    }

    childCount := a.traverser.ast.GetChildCount(node.Type)

    if childCount == VariableChildren {
        for !a.traverser.IsAtEnd() {
            node := a.traverser.Next()

            if node.Type == EndNode {
                return true
            }

            if !a.displayNode(node, depth + 1) {
                return false
            }
        }

        return a.tryWrite(
            "%sMissing EndNode\n",
            strings.Repeat(" ", int(spacesPerLevel * (depth + 1))),
        )
    }

    for i := range childCount {
        if a.traverser.IsAtEnd() {
            return a.tryWrite(
                "%s%d missing children\n",
                strings.Repeat(" ", int(spacesPerLevel * (depth + 1))),
                childCount - i,
            )
        }

        node := a.traverser.Next()

        if !a.displayNode(node, depth + 1) {
            return false
        }
    }

    return true
}

func (a *astDebugger) tryWrite(format string, args ...any) bool {
    _, err := fmt.Fprintf(a.output, format, args...)

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
