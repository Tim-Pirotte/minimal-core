package logrendering

import (
	"bytes"
	messaging "minimal/minimal-core/built-in/messaging"
	"sort"
	"strconv"
	"strings"
)

func (l *LogRenderer) outputContext(bb *bytes.Buffer, ctx messaging.CodeContext) {
    if bb == nil {
        panic("an uninitialized bytes buffer was passed to the context renderer")
    }

    if ctx.StartLineNumber < uint(len(ctx.LinesBefore)) {
        panic("lines before the lines in focus can be negative")
    }

    // Header
    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "source_start", "["))
    bb.WriteString(resetAnsi)

    bb.WriteString(l.getStrOrDefault("context", "source_color", ""))
    bb.WriteString(ctx.Source)
    bb.WriteString(resetAnsi)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "source_line_separator", ":"))
    bb.WriteString(resetAnsi)

    bb.WriteString(l.getStrOrDefault("context", "start_line_color", ""))
    bb.WriteString(strconv.FormatUint(uint64(ctx.StartLineNumber), 10))
    bb.WriteString(resetAnsi)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "source_end", "]"))
    bb.WriteString(resetAnsi)

    bb.WriteString(strings.Repeat("\n", l.getIntOrDefault("context", "source_context_padding", 1) + 1))

    largestAmountOfDigits := countDigits(ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) - 1 + uint(len(ctx.LinesAfter)))

    leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

    // Window top
    bb.WriteString(leftPadding)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "context_window_top_symbol", ""))
    bb.WriteString(resetAnsi)

    bb.WriteString("\n")

    // Lines out of focus before
    for i, line := range ctx.LinesBefore {
        lineNumber := ctx.StartLineNumber - uint(len(ctx.LinesBefore)) + uint(i)

        l.renderLinePrefix(
            bb,
            lineNumber,
            largestAmountOfDigits,
            l.getStrOrDefault("context", "out_of_focus_line_number_color", ""),
        )

        bb.WriteString(line)
        bb.WriteString("\n")
    }

    // Lines in focus
    for i, line := range ctx.LinesInFocus {
        lineNumber := ctx.StartLineNumber + uint(i)

        l.renderLinePrefix(bb, lineNumber, largestAmountOfDigits, l.getStrOrDefault("context", "in_focus_line_number_color", ""))

        bb.WriteString(line.Content)
        bb.WriteString("\n")

        sort.Slice(line.Annotations, func(i, j int) bool {
            return line.Annotations[i].Range.Start < line.Annotations[j].Range.Start
        })

        // Annotation lines
        currentLineEnd := uint(0)

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(bb, leftPadding, l.getStrOrDefault("context", "annotation_line_start_symbol", "-"))
        }

        for _, annotation := range line.Annotations {
            if annotation.Range.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(bb, leftPadding, l.getStrOrDefault("context", "annotation_line_start_symbol", "-"))
                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))

            bb.WriteString(l.getSeverityColor(annotation.Severity))
            bb.WriteString(strings.Repeat(l.getStrOrDefault("context", "annotation_symbol", "_"), int(annotation.Range.Length)))
            bb.WriteString(resetAnsi)

            currentLineEnd = annotation.Range.Start + annotation.Range.Length
        }

        bb.WriteString("\n")

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(bb, leftPadding, l.getStrOrDefault("context", "annotation_comment_start_symbol", ":"))
        }

        // Annotation comments
        currentLineEnd = 0

        for _, annotation := range line.Annotations {
            if annotation.Range.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(bb, leftPadding, l.getStrOrDefault("context", "annotation_comment_start_symbol", ":"))
                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))

            if l.getBoolOrDefault("context", "text_has_severity_color", true) {
                bb.WriteString(l.getSeverityColor(annotation.Severity))
            }

            bb.WriteString(annotation.Message)

            if l.getBoolOrDefault("context", "text_has_severity_color", true) {
                bb.WriteString(resetAnsi)
            }

            currentLineEnd = annotation.Range.Start + annotation.Range.Length
        }

        bb.WriteString("\n")
    }

    // Lines after
    for i, line := range ctx.LinesAfter {
        lineNumber := ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) + uint(i)

        l.renderLinePrefix(bb, lineNumber, largestAmountOfDigits, l.getStrOrDefault("context", "out_of_focus_line_number_color", ""))
        bb.WriteString(line)
        bb.WriteString("\n")
    }

    // Window bottom
    bb.WriteString(leftPadding)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "context_window_bottom_symbol", ""))
    bb.WriteString(resetAnsi)

    bb.WriteString("\n")
}

func (l *LogRenderer) renderAnnotationPrefix(bb *bytes.Buffer, leftPadding, prefixSymbol string) {
    bb.WriteString(leftPadding)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(prefixSymbol)
    bb.WriteString(resetAnsi)

    bb.WriteString(" ")
}

func (l *LogRenderer) renderLinePrefix(bb *bytes.Buffer, lineNumber uint, largestAmountOfDigits int, color string) {
    renderLineNumber(bb, lineNumber, largestAmountOfDigits, color)

    bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
    bb.WriteString(l.getStrOrDefault("context", "line_count_sep", "|"))
    bb.WriteString(resetAnsi)

    bb.WriteString(" ")
}
