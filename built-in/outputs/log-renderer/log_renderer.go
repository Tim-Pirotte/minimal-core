package logrendering

import (
	"bytes"
	"fmt"
	"io"
	"minimal/minimal-core/built-in/messenger"
)

type LogRenderer struct {
    writer io.Writer
    Config Config
}

type Config struct {
    ResetAnsi     string
    SymbolColor   string
    Severity      SeverityConfig
    Context       ContextConfig
    Diff          DiffConfig
    InfoReference InfoReferenceConfig
}

type SeverityConfig struct {
    VerboseColor       string
    DebugColor         string
    InfoColor          string
    WarningColor       string
    SevereWarningColor string
    ErrorColor         string
    CriticalColor      string
}

type ContextConfig struct {
    SourceContextPadding      int
    SourceStart               string
    SourceEnd                 string
    SourceLineSeparator       string
    WindowTop                 string
    WindowBottom              string
    LineCountSeparator        string
    SourceColor               string
    StartLineColor            string
    Annotation                string
    AnnotationLineStart       string
    AnnotationCommentStart    string
    OutOfFocusLineNumberColor string
    InFocusLineNumberColor    string
}

type DiffConfig struct {
    WindowTop                 string
    WindowBottom              string
    LineCountSeparator        string
    OutOfFocusLineNumberColor string
    InFocusLineNumberColor    string
    RemoveLinePrefix          string
    RemoveLineColor           string
    AddLinePrefix             string
    AddLineColor              string
}

type InfoReferenceConfig struct {
    HintPrefixColor string
    HintColor       string
    MoreInfoColor   string
}

func NewLogRenderer(writer io.Writer) *LogRenderer {
    return &LogRenderer{
        writer,
        Config{
            ResetAnsi: "\033[0m",
            SymbolColor: "\u001b[38;2;255;255;255m",
            Severity: SeverityConfig{
                VerboseColor: "\u001b[38;2;115;115;115m",
                DebugColor: "\u001b[38;2;75;177;229m",
                InfoColor: "\u001b[38;2;177;73;230m",
                WarningColor: "\u001b[38;2;236;236;73m",
                SevereWarningColor: "\u001b[38;2;246;152;58m",
                ErrorColor: "\u001b[38;2;234;46;46m",
                CriticalColor: "\u001b[38;2;163;21;21m",
            },
            Context: ContextConfig{
                SourceContextPadding: 1,
                SourceStart: "[",
                SourceEnd: "]",
                SourceLineSeparator: ":",
                WindowTop: "╭──",
                WindowBottom: "╰──",
                LineCountSeparator: "│",
                SourceColor: "\u001b[38;2;93;209;93m",
                StartLineColor: "\u001b[38;2;242;205;205m",
                Annotation: "─",
                AnnotationLineStart: "•",
                AnnotationCommentStart: "∘",
                OutOfFocusLineNumberColor: "\u001b[38;2;115;115;115m",
                InFocusLineNumberColor: "\u001b[38;2;75;177;229m",
            },
            Diff: DiffConfig{
                WindowTop: "╭──",
                WindowBottom: "╰──",
                OutOfFocusLineNumberColor: "\u001b[38;2;115;115;115m",
                InFocusLineNumberColor: "\u001b[38;2;75;177;229m",
                LineCountSeparator: "│",
                RemoveLinePrefix: "-",
                RemoveLineColor: "\u001b[38;2;234;46;46m",
                AddLinePrefix: "+",
                AddLineColor: "\u001b[38;2;93;209;93m",
            },
            InfoReference: InfoReferenceConfig{
                HintPrefixColor: "\u001b[38;2;75;0;229m",
                HintColor: "\u001b[38;2;75;177;229m",
                MoreInfoColor: "\u001b[38;2;75;0;229m",
            },
        },
    }
}

func (l *LogRenderer) Receive(m messenger.Message) {
    bb := bytes.NewBuffer(make([]byte, 0))

	bb.WriteString(l.Config.ResetAnsi)

    bb.WriteString(l.getSeverityColor(m.Severity))
	fmt.Fprintf(bb, "%14s", stringifySeverity(m.Severity))
	bb.WriteString(l.Config.ResetAnsi)
	bb.WriteString(l.Config.SymbolColor)
    bb.WriteString(": ") // TODO make configurable
	bb.WriteString(l.Config.ResetAnsi)
	bb.WriteString(m.Message)
	bb.WriteString("\n")

    // TODO add proper output
    for _, span := range m.Context {
        bb.WriteString("\"")
        bb.WriteString(span.Content)
        bb.WriteString("\"")

        if span.Note != "" {
            bb.WriteString(" (")
            bb.WriteString(span.Note)
            bb.WriteString(")")
        }

        bb.WriteString("\n")
    }

    for _, span := range m.AdditionalContext {
        bb.WriteString("Additional: ")
        bb.WriteString("\"")
        bb.WriteString(span.Content)
        bb.WriteString("\"")

        if span.Note != "" {
            bb.WriteString(" (")
            bb.WriteString(span.Note)
            bb.WriteString(")")
        }

        bb.WriteString("\n")
    }

    for _, suggestion := range m.Suggestions {
        bb.WriteString("Suggestion: ")
        bb.WriteString(suggestion.Suggestion)

        for _, replacement := range suggestion.Replacements {
            bb.WriteString("- ")
            bb.WriteString("\"")
            bb.WriteString(replacement.From.Content)
            bb.WriteString("\"")

            if replacement.From.Note != "" {
                bb.WriteString(" (")
                bb.WriteString(replacement.From.Note)
                bb.WriteString(")")
            }

            bb.WriteString(" ")

            bb.WriteString("+ ")
            bb.WriteString("\"")
            bb.WriteString(replacement.To.Content)
            bb.WriteString("\"")

            if replacement.To.Note != "" {
                bb.WriteString(" (")
                bb.WriteString(replacement.To.Note)
                bb.WriteString(")")
            }

            bb.WriteString("\n")
        }
    }

    for _, note := range m.Notes {
        bb.WriteString("Note: ")
        bb.WriteString(note)
        bb.WriteString("\n")
    }

    bb.WriteString("\n")

    _, err := l.writer.Write(bb.Bytes())

    if err != nil {
        // TODO what here?
    }
}

func (c *Config) RemoveANSI() {
    c.ResetAnsi = ""
    c.SymbolColor = ""

    c.Severity.VerboseColor = ""
    c.Severity.DebugColor = ""
    c.Severity.InfoColor = ""
    c.Severity.WarningColor = ""
    c.Severity.SevereWarningColor = ""
    c.Severity.ErrorColor = ""
    c.Severity.CriticalColor = ""

    c.Context.SourceColor = ""
    c.Context.StartLineColor = ""
    c.Context.OutOfFocusLineNumberColor = ""
    c.Context.InFocusLineNumberColor = ""

    c.Diff.OutOfFocusLineNumberColor = ""
    c.Diff.InFocusLineNumberColor = ""
    c.Diff.RemoveLineColor = ""
    c.Diff.AddLineColor = ""

    c.InfoReference.HintPrefixColor = ""
    c.InfoReference.HintColor = ""
    c.InfoReference.MoreInfoColor = ""
}

func (c *Config) RemoveUnicode() {
    c.Context.WindowTop = ""
    c.Context.WindowBottom = ""
    c.Context.LineCountSeparator = "|"
    c.Context.Annotation = "-"
    c.Context.AnnotationLineStart = "*"
    c.Context.AnnotationCommentStart = "o"

    c.Diff.WindowTop = ""
    c.Diff.WindowBottom = ""
    c.Diff.LineCountSeparator = "|"
}

func stringifySeverity(s messenger.Severity) string {
	var sAsStr string

	switch s {
	case messenger.Verbose:
		sAsStr = "VERBOSE"
	case messenger.Debug:
		sAsStr = "DEBUG"
	case messenger.Info:
		sAsStr = "INFO"
	case messenger.Warning:
		sAsStr = "WARNING"
	case messenger.SevereWarning:
		sAsStr = "SEVERE WARNING"
	case messenger.Error:
		sAsStr = "ERROR"
	case messenger.Critical:
		sAsStr = "CRITICAL"
	default:
		panic(fmt.Sprintf("missing string representation for the enum Severity: %d", s))
	}

	return sAsStr
}

func (l *LogRenderer) getSeverityColor(s messenger.Severity) string {
	var sAsStr string

	switch s {
	case messenger.Verbose:
		sAsStr = l.Config.Severity.VerboseColor
	case messenger.Debug:
		sAsStr = l.Config.Severity.DebugColor
	case messenger.Info:
		sAsStr = l.Config.Severity.InfoColor
	case messenger.Warning:
		sAsStr = l.Config.Severity.WarningColor
	case messenger.SevereWarning:
		sAsStr = l.Config.Severity.SevereWarningColor
	case messenger.Error:
		sAsStr = l.Config.Severity.ErrorColor
	case messenger.Critical:
		sAsStr = l.Config.Severity.CriticalColor
	default:
		panic(fmt.Sprintf("missing color for the enum Severity: %d", s))
	}

	return sAsStr
}

// func (l *LogRenderer) renderHint(bb *bytes.Buffer, hint messaging.Hint) {
// 	if hint.Text != "" {
// 		bb.WriteString(l.Config.InfoReference.HintPrefixColor)
// 		bb.WriteString("Hint")
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(l.Config.SymbolColor)
// 		bb.WriteString(": ") // TODO Make configurable
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(l.Config.InfoReference.HintColor)
// 		bb.WriteString(hint.Text)
// 		bb.WriteString(l.Config.ResetAnsi)
// 		bb.WriteString("\n")
// 	}

// 	if hint.MoreInfoReference != "" {
// 		bb.WriteString(l.Config.InfoReference.MoreInfoColor)
// 		bb.WriteString("More info about this")
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(l.Config.SymbolColor)
// 		bb.WriteString(": ")
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(hint.MoreInfoReference)
// 		bb.WriteString("\n")
// 	}
// }

// func (l *LogRenderer) renderContext(bb *bytes.Buffer, ctx messaging.CodeContext) {
//     if bb == nil {
//         panic("an uninitialized bytes buffer was passed to the context renderer")
//     }

//     if ctx.StartLineNumber < uint(len(ctx.LinesBefore)) {
//         panic("lines before the lines in focus can be negative")
//     }

//     // Header
//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.SourceStart)
//     bb.WriteString(l.Config.ResetAnsi)

//     bb.WriteString(l.Config.Context.SourceColor)
//     bb.WriteString(ctx.Source)
//     bb.WriteString(l.Config.ResetAnsi)

//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.SourceLineSeparator)
//     bb.WriteString(l.Config.ResetAnsi)

//     bb.WriteString(l.Config.Context.StartLineColor)
//     bb.WriteString(strconv.FormatUint(uint64(ctx.StartLineNumber), 10))
//     bb.WriteString(l.Config.ResetAnsi)

//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.SourceEnd)
//     bb.WriteString(l.Config.ResetAnsi)
//     bb.WriteString(strings.Repeat("\n", l.Config.Context.SourceContextPadding + 1))

//     largestAmountOfDigits := countDigits(ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) - 1 + uint(len(ctx.LinesAfter)))
//     leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

//     // Window top
//     bb.WriteString(leftPadding)
//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.WindowTop)
//     bb.WriteString(l.Config.ResetAnsi)
//     bb.WriteString("\n")

//     // Lines out of focus before
//     for i, line := range ctx.LinesBefore {
//         lineNumber := ctx.StartLineNumber - uint(len(ctx.LinesBefore)) + uint(i)

//         l.renderLinePrefix(
//             bb,
//             lineNumber,
//             largestAmountOfDigits,
//             l.Config.Context.OutOfFocusLineNumberColor,
//         )

//         bb.WriteString(line)
//         bb.WriteString("\n")
//     }

//     // Lines in focus
//     for i, line := range ctx.LinesInFocus {
//         lineNumber := ctx.StartLineNumber + uint(i)

//         l.renderLinePrefix(
//             bb,
//             lineNumber,
//             largestAmountOfDigits,
//             l.Config.Context.InFocusLineNumberColor,
//         )

//         bb.WriteString(line.Content)
//         bb.WriteString("\n")

//         sort.Slice(line.Annotations, func(i, j int) bool {
//             return line.Annotations[i].Range.Start < line.Annotations[j].Range.Start
//         })

//         // Annotation lines
//         currentLineEnd := uint(0)

//         if len(line.Annotations) > 0 {
//             l.renderAnnotationPrefix(
//                 bb,
//                 leftPadding,
//                 l.Config.Context.AnnotationLineStart,
//             )
//         }

//         for _, annotation := range line.Annotations {
//             if annotation.Range.Start < currentLineEnd {
//                 bb.WriteString("\n")
//                 l.renderAnnotationPrefix(
//                     bb,
//                     leftPadding,
//                     l.Config.Context.AnnotationLineStart,
//                 )

//                 currentLineEnd = 0
//             }

//             bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))
//             bb.WriteString(l.getSeverityColor(annotation.Severity))
//             bb.WriteString(strings.Repeat(l.Config.Context.Annotation, int(annotation.Range.Length)))
//             bb.WriteString(l.Config.ResetAnsi)

//             currentLineEnd = annotation.Range.Start + annotation.Range.Length
//         }

//         bb.WriteString("\n")

//         if len(line.Annotations) > 0 {
//             l.renderAnnotationPrefix(
//                 bb,
//                 leftPadding,
//                 l.Config.Context.AnnotationCommentStart,
//             )
//         }

//         // Annotation comments
//         currentLineEnd = 0

//         for _, annotation := range line.Annotations {
//             if annotation.Range.Start < currentLineEnd {
//                 bb.WriteString("\n")
//                 l.renderAnnotationPrefix(
//                     bb,
//                     leftPadding,
//                     l.Config.Context.AnnotationCommentStart,
//                 )

//                 currentLineEnd = 0
//             }

//             bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))
//             bb.WriteString(annotation.Message)

//             currentLineEnd = annotation.Range.Start + annotation.Range.Length
//         }

//         bb.WriteString("\n")
//     }

//     // Lines after
//     for i, line := range ctx.LinesAfter {
//         lineNumber := ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) + uint(i)

//         l.renderLinePrefix(
//             bb,
//             lineNumber,
//             largestAmountOfDigits,
//             l.Config.Context.OutOfFocusLineNumberColor,
//         )

//         bb.WriteString(line)
//         bb.WriteString("\n")
//     }

//     // Window bottom
//     bb.WriteString(leftPadding)
//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.WindowBottom)
//     bb.WriteString(l.Config.ResetAnsi)
//     bb.WriteString("\n")
// }

// func (l *LogRenderer) renderAnnotationPrefix(bb *bytes.Buffer, leftPadding, prefixSymbol string) {
//     bb.WriteString(leftPadding)
//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(prefixSymbol)
//     bb.WriteString(l.Config.ResetAnsi)
//     bb.WriteString(" ")
// }

// func (l *LogRenderer) renderLinePrefix(bb *bytes.Buffer, lineNumber uint, largestAmountOfDigits int, color string) {
//     l.renderLineNumber(bb, lineNumber, largestAmountOfDigits, color)
//     bb.WriteString(l.Config.SymbolColor)
//     bb.WriteString(l.Config.Context.LineCountSeparator)
//     bb.WriteString(l.Config.ResetAnsi)
//     bb.WriteString(" ")
// }

// func (l *LogRenderer) renderDiff(bb *bytes.Buffer, diff messaging.Diff) {
// 	if bb == nil {
//         panic("an uninitialized bytes buffer was passed to the diff renderer")
//     }

// 	if diff.StartLineNumber < uint(len(diff.LinesBefore)) {
//         panic("lines before the lines in focus can be negative")
//     }

// 	largestAmountOfDigits := countDigits(
// 		diff.StartLineNumber + uint(max(len(diff.LinesToAdd) + len(diff.LinesToRemove))) - 1 + uint(len(diff.LinesAfter)),
// 	)

// 	leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

// 	// Window top
// 	bb.WriteString(leftPadding)

// 	bb.WriteString(l.Config.SymbolColor)
// 	bb.WriteString(l.Config.Diff.WindowTop)
// 	bb.WriteString(l.Config.ResetAnsi)

// 	bb.WriteString("\n")

// 	// Lines out of focus before
// 	for i, line := range diff.LinesBefore {
// 		l.renderLineNumber(
// 			bb,
// 			diff.StartLineNumber - uint(len(diff.LinesBefore)) + uint(i),
// 			largestAmountOfDigits,
// 			l.Config.Diff.OutOfFocusLineNumberColor,
// 		)

// 		bb.WriteString(l.Config.SymbolColor)
// 		bb.WriteString(l.Config.Diff.LineCountSeparator)
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(" ")
// 		bb.WriteString(line)
// 		bb.WriteString("\n")
// 	}

// 	// Lines to remove
// 	for i, line := range diff.LinesToRemove {
// 		l.renderLineNumber(
// 			bb,
// 			diff.StartLineNumber + uint(i),
// 			largestAmountOfDigits,
// 			l.Config.Diff.InFocusLineNumberColor,
// 		)

// 		bb.WriteString(l.Config.Diff.RemoveLineColor)
// 		bb.WriteString(l.Config.Diff.RemoveLinePrefix)
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(" ")
// 		bb.WriteString(line)
// 		bb.WriteString("\n")
// 	}

// 	// Lines to add
// 	for i, line := range diff.LinesToAdd {
// 		l.renderLineNumber(
// 			bb,
// 			diff.StartLineNumber + uint(i),
// 			largestAmountOfDigits,
// 			l.Config.Diff.InFocusLineNumberColor,
// 		)

// 		bb.WriteString(l.Config.Diff.AddLineColor)
// 		bb.WriteString(l.Config.Diff.AddLinePrefix)
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(" ")
// 		bb.WriteString(line)
// 		bb.WriteString("\n")
// 	}

// 	// Lines out of focus after
// 	for i, line := range diff.LinesAfter {
// 		l.renderLineNumber(
// 			bb,
// 			diff.StartLineNumber + uint(len(diff.LinesToAdd)) + uint(i),
// 			largestAmountOfDigits,
// 			l.Config.Diff.OutOfFocusLineNumberColor,
// 		)

// 		bb.WriteString(l.Config.SymbolColor)
// 		bb.WriteString(l.Config.Diff.LineCountSeparator)
// 		bb.WriteString(l.Config.ResetAnsi)

// 		bb.WriteString(" ")
// 		bb.WriteString(line)
// 		bb.WriteString("\n")
// 	}

// 	// Window bottom
// 	bb.WriteString(leftPadding)

// 	bb.WriteString(l.Config.SymbolColor)
// 	bb.WriteString(l.Config.Diff.WindowBottom)
// 	bb.WriteString(l.Config.ResetAnsi)

// 	bb.WriteString("\n")
// }

// func countDigits(number uint) int {
// 	if number == 0 {
// 		return 1
// 	}

// 	count := 0

// 	for number > 0 {
// 		number /= 10
// 		count++
// 	}

// 	return count
// }

// func (l *LogRenderer) renderLineNumber(
// 	bb *bytes.Buffer,
// 	lineNumber uint,
// 	largestAmountOfDigits int,
// 	color string,
// ) {
// 	lineNumberAsStr := strconv.FormatUint(uint64(lineNumber), 10)

// 	bb.WriteString(strings.Repeat(" ", largestAmountOfDigits-len(lineNumberAsStr)))
// 	bb.WriteString(color)
// 	bb.WriteString(lineNumberAsStr)
// 	bb.WriteString(l.Config.ResetAnsi)
// 	bb.WriteString(" ")
// }
