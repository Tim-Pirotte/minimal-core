package formattedtext

import (
	"bytes"
	"sort"
	"strconv"
	"strings"
)

// A part of the source code with highlighted lines and annotations
func (l *LogRenderer) RenderLogContext(bb *bytes.Buffer, ctx CodeContext) {
    if bb == nil {
        panic("an uninitialized bytes buffer was passed to the context renderer")
    }

    if ctx.StartLineNumber < uint(len(ctx.LinesBefore)) {
        panic("lines before the lines in focus can be negative")
    }

    // Header
    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.SourceStart)
    bb.WriteString(resetAnsi)

    bb.WriteString(l.config.Context.SourceColor)
    bb.WriteString(ctx.Source)
    bb.WriteString(resetAnsi)

    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.SourceLineSep)
    bb.WriteString(resetAnsi)
    
    bb.WriteString(l.config.Context.StartLineColor)
    bb.WriteString(strconv.FormatUint(uint64(ctx.StartLineNumber), 10))
    bb.WriteString(resetAnsi)
    
    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.SourceEnd)
    bb.WriteString(resetAnsi)
    
    bb.WriteString(strings.Repeat("\n", l.config.Context.SourceContextPadding + 1))

    largestAmountOfDigits := countDigits(ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) - 1 + uint(len(ctx.LinesAfter)))

    leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

    // Window top
    bb.WriteString(leftPadding)

    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.ContextWindowTopSymbol)
    bb.WriteString(resetAnsi)

    bb.WriteString("\n")
    
    // Lines out of focus before
    for i, line := range ctx.LinesBefore {
        lineNumber := ctx.StartLineNumber - uint(len(ctx.LinesBefore)) + uint(i)

        l.renderLinePrefix(bb, lineNumber, largestAmountOfDigits, l.config.Context.OutOfFocusLineNumberColor)
        
        bb.WriteString(line)
        bb.WriteString("\n")
    }
    
    // Lines in focus
    for i, line := range ctx.LinesInFocus {
        lineNumber := ctx.StartLineNumber + uint(i)

        l.renderLinePrefix(bb, lineNumber, largestAmountOfDigits, l.config.Context.InFocusLineNumberColor)
        
        bb.WriteString(line.Content)
        bb.WriteString("\n")

        sort.Slice(line.Annotations, func(i, j int) bool {
            return line.Annotations[i].Start < line.Annotations[j].Start
        })

        // Annotation lines
        currentLineEnd := 0

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(bb, leftPadding, l.config.Context.AnnotationLineStartSymbol)
        }

        for _, annotation := range line.Annotations {
            if annotation.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(bb, leftPadding, l.config.Context.AnnotationLineStartSymbol)
                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", annotation.Start - currentLineEnd))

            bb.WriteString(l.config.Severity.color(annotation.Severity))
            bb.WriteString(strings.Repeat(l.config.Context.AnnotationSymbol, annotation.Length))
            bb.WriteString(resetAnsi)

            currentLineEnd = annotation.Start + annotation.Length
        }

        bb.WriteString("\n")

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(bb, leftPadding, l.config.Context.AnnotationCommentStartSymbol)
        }

        // Annotation comments
        currentLineEnd = 0

        for _, annotation := range line.Annotations {
            if annotation.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(bb, leftPadding, l.config.Context.AnnotationCommentStartSymbol)
                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", annotation.Start - currentLineEnd))
            
            if l.config.Context.TextHasSeverityColor {
                bb.WriteString(l.config.Severity.color(annotation.Severity))
            }
            
            bb.WriteString(annotation.Message)

            if l.config.Context.TextHasSeverityColor {
                bb.WriteString(resetAnsi)
            }
            
            currentLineEnd = annotation.Start + annotation.Length
        }

        bb.WriteString("\n")
    }
    
    // Lines after
    for i, line := range ctx.LinesAfter {
        lineNumber := ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) + uint(i)

        l.renderLinePrefix(bb, lineNumber, largestAmountOfDigits, l.config.Context.OutOfFocusLineNumberColor)
        bb.WriteString(line)
        bb.WriteString("\n")
    }
    
    // Window bottom
    bb.WriteString(leftPadding)
    
    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.ContextWindowBottomSymbol)
    bb.WriteString(resetAnsi)
    
    bb.WriteString("\n")
}

func (l *LogRenderer) renderAnnotationPrefix(bb *bytes.Buffer, leftPadding, prefixSymbol string) {
    bb.WriteString(leftPadding)
            
    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(prefixSymbol)
    bb.WriteString(resetAnsi)
    
    bb.WriteString(" ")
}

func (l *LogRenderer) renderLinePrefix(bb *bytes.Buffer, lineNumber uint, largestAmountOfDigits int, color string) {
    renderLineNumber(bb, lineNumber, largestAmountOfDigits, color)
    
    bb.WriteString(l.config.General.SymbolColor)
    bb.WriteString(l.config.Context.LineCountSep)
    bb.WriteString(resetAnsi)
    
    bb.WriteString(" ")
}

type Annotation struct {
	Start    int
	Length   int
	Message  string
	Severity Severity
}

type Line struct {
	Content     string
	Annotations []Annotation
}

type CodeContext struct {
    Source          string
	StartLineNumber uint
	LinesBefore     []string
	LinesInFocus    []Line
	LinesAfter      []string
}
