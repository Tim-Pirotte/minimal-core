package logrendering

import (
	"bytes"
	"fmt"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/messaging"
	"sort"
	"strconv"
	"strings"
)

type LogRenderer struct {
    logger logging.Logger
    writer io.Writer
    config Config
}

type Config struct {
    SymbolColor   string
    ResetAnsi     string
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

func NewLogRenderer(sourceGen *logging.SourceGenerator, writer io.Writer, config Config) *LogRenderer {
    logger, _ := sourceGen.GetLogger("logRendering")

    return &LogRenderer{logger, writer, config}
}

func (l *LogRenderer) Receive(messageParts []messaging.MessagePart) {
    bb := bytes.NewBuffer(make([]byte, 0))

    for _, part := range messageParts {
        switch p := part.(type) {
        case *messaging.Message:
            l.renderMessage(bb, *p)
        case *messaging.CodeContext:
            l.renderContext(bb, *p)
        case *messaging.Hint:
            l.renderHint(bb, *p)
        case *messaging.Diff:
            l.renderDiff(bb, *p)
        default:
            // TODO change to proper error
            panic("unsuported message part")
        }
    }

    l.writer.Write([]byte(l.config.ResetAnsi))
    _, err := l.writer.Write(bb.Bytes())

    if err != nil {
        l.logger.Error().Err(err).Msg("unsuccessfull write")
    }
}

func stringifySeverity(s messaging.Severity) string {
	var sAsStr string

	switch s {
	case messaging.Verbose:
		sAsStr = "VERBOSE"
	case messaging.Debug:
		sAsStr = "DEBUG"
	case messaging.Info:
		sAsStr = "INFO"
	case messaging.Warning:
		sAsStr = "WARNING"
	case messaging.SevereWarning:
		sAsStr = "SEVERE WARNING"
	case messaging.Error:
		sAsStr = "ERROR"
	case messaging.Critical:
		sAsStr = "CRITICAL"
	default:
		panic(fmt.Sprintf("missing string representation for the enum Severity: %d", s))
	}

	return sAsStr
}

func (l *LogRenderer) getSeverityColor(s messaging.Severity) string {
	var sAsStr string

	switch s {
	case messaging.Verbose:
		sAsStr = l.config.Severity.VerboseColor
	case messaging.Debug:
		sAsStr = l.config.Severity.DebugColor
	case messaging.Info:
		sAsStr = l.config.Severity.InfoColor
	case messaging.Warning:
		sAsStr = l.config.Severity.WarningColor
	case messaging.SevereWarning:
		sAsStr = l.config.Severity.SevereWarningColor
	case messaging.Error:
		sAsStr = l.config.Severity.ErrorColor
	case messaging.Critical:
		sAsStr = l.config.Severity.CriticalColor
	default:
		panic(fmt.Sprintf("missing color for the enum Severity: %d", s))
	}

	return sAsStr
}

func (l *LogRenderer) renderMessage(bb *bytes.Buffer, m messaging.Message) {
	bb.WriteString(l.getSeverityColor(m.Severity))
	bb.WriteString(stringifySeverity(m.Severity))
	bb.WriteString(l.config.ResetAnsi)
	bb.WriteString(" ")
	bb.WriteString(m.Category)
	bb.WriteString(l.config.SymbolColor)
	bb.WriteString(": ") // TODO make configurable
	bb.WriteString(l.config.ResetAnsi)
	bb.WriteString(m.Message)
	bb.WriteString("\n")
}

func (l *LogRenderer) renderHint(bb *bytes.Buffer, hint messaging.Hint) {
	if hint.Text != "" {
		bb.WriteString(l.config.InfoReference.HintPrefixColor)
		bb.WriteString("Hint")
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(l.config.SymbolColor)
		bb.WriteString(": ") // TODO Make configurable
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(l.config.InfoReference.HintColor)
		bb.WriteString(hint.Text)
		bb.WriteString(l.config.ResetAnsi)
		bb.WriteString("\n")
	}

	if hint.MoreInfoReference != "" {
		bb.WriteString(l.config.InfoReference.MoreInfoColor)
		bb.WriteString("More info about this")
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(l.config.SymbolColor)
		bb.WriteString(": ")
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(hint.MoreInfoReference)
		bb.WriteString("\n")
	}
}

func (l *LogRenderer) renderContext(bb *bytes.Buffer, ctx messaging.CodeContext) {
    if bb == nil {
        panic("an uninitialized bytes buffer was passed to the context renderer")
    }

    if ctx.StartLineNumber < uint(len(ctx.LinesBefore)) {
        panic("lines before the lines in focus can be negative")
    }

    // Header
    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.SourceStart)
    bb.WriteString(l.config.ResetAnsi)

    bb.WriteString(l.config.Context.SourceColor)
    bb.WriteString(ctx.Source)
    bb.WriteString(l.config.ResetAnsi)

    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.SourceLineSeparator)
    bb.WriteString(l.config.ResetAnsi)

    bb.WriteString(l.config.Context.StartLineColor)
    bb.WriteString(strconv.FormatUint(uint64(ctx.StartLineNumber), 10))
    bb.WriteString(l.config.ResetAnsi)

    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.SourceEnd)
    bb.WriteString(l.config.ResetAnsi)
    bb.WriteString(strings.Repeat("\n", l.config.Context.SourceContextPadding + 1))

    largestAmountOfDigits := countDigits(ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) - 1 + uint(len(ctx.LinesAfter)))
    leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

    // Window top
    bb.WriteString(leftPadding)
    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.WindowTop)
    bb.WriteString(l.config.ResetAnsi)
    bb.WriteString("\n")

    // Lines out of focus before
    for i, line := range ctx.LinesBefore {
        lineNumber := ctx.StartLineNumber - uint(len(ctx.LinesBefore)) + uint(i)

        l.renderLinePrefix(
            bb,
            lineNumber,
            largestAmountOfDigits,
            l.config.Context.OutOfFocusLineNumberColor,
        )

        bb.WriteString(line)
        bb.WriteString("\n")
    }

    // Lines in focus
    for i, line := range ctx.LinesInFocus {
        lineNumber := ctx.StartLineNumber + uint(i)

        l.renderLinePrefix(
            bb,
            lineNumber,
            largestAmountOfDigits,
            l.config.Context.InFocusLineNumberColor,
        )

        bb.WriteString(line.Content)
        bb.WriteString("\n")

        sort.Slice(line.Annotations, func(i, j int) bool {
            return line.Annotations[i].Range.Start < line.Annotations[j].Range.Start
        })

        // Annotation lines
        currentLineEnd := uint(0)

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(
                bb,
                leftPadding,
                l.config.Context.AnnotationLineStart,
            )
        }

        for _, annotation := range line.Annotations {
            if annotation.Range.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(
                    bb,
                    leftPadding,
                    l.config.Context.AnnotationLineStart,
                )

                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))
            bb.WriteString(l.getSeverityColor(annotation.Severity))
            bb.WriteString(strings.Repeat(l.config.Context.Annotation, int(annotation.Range.Length)))
            bb.WriteString(l.config.ResetAnsi)

            currentLineEnd = annotation.Range.Start + annotation.Range.Length
        }

        bb.WriteString("\n")

        if len(line.Annotations) > 0 {
            l.renderAnnotationPrefix(
                bb,
                leftPadding,
                l.config.Context.AnnotationCommentStart,
            )
        }

        // Annotation comments
        currentLineEnd = 0

        for _, annotation := range line.Annotations {
            if annotation.Range.Start < currentLineEnd {
                bb.WriteString("\n")
                l.renderAnnotationPrefix(
                    bb,
                    leftPadding,
                    l.config.Context.AnnotationCommentStart,
                )

                currentLineEnd = 0
            }

            bb.WriteString(strings.Repeat(" ", int(annotation.Range.Start - currentLineEnd)))
            bb.WriteString(annotation.Message)

            currentLineEnd = annotation.Range.Start + annotation.Range.Length
        }

        bb.WriteString("\n")
    }

    // Lines after
    for i, line := range ctx.LinesAfter {
        lineNumber := ctx.StartLineNumber + uint(len(ctx.LinesInFocus)) + uint(i)

        l.renderLinePrefix(
            bb,
            lineNumber,
            largestAmountOfDigits,
            l.config.Context.OutOfFocusLineNumberColor,
        )

        bb.WriteString(line)
        bb.WriteString("\n")
    }

    // Window bottom
    bb.WriteString(leftPadding)
    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.WindowBottom)
    bb.WriteString(l.config.ResetAnsi)
    bb.WriteString("\n")
}

func (l *LogRenderer) renderAnnotationPrefix(bb *bytes.Buffer, leftPadding, prefixSymbol string) {
    bb.WriteString(leftPadding)
    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(prefixSymbol)
    bb.WriteString(l.config.ResetAnsi)
    bb.WriteString(" ")
}

func (l *LogRenderer) renderLinePrefix(bb *bytes.Buffer, lineNumber uint, largestAmountOfDigits int, color string) {
    l.renderLineNumber(bb, lineNumber, largestAmountOfDigits, color)
    bb.WriteString(l.config.SymbolColor)
    bb.WriteString(l.config.Context.LineCountSeparator)
    bb.WriteString(l.config.ResetAnsi)
    bb.WriteString(" ")
}

func (l *LogRenderer) renderDiff(bb *bytes.Buffer, diff messaging.Diff) {
	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the diff renderer")
    }

	if diff.StartLineNumber < uint(len(diff.LinesBefore)) {
        panic("lines before the lines in focus can be negative")
    }

	largestAmountOfDigits := countDigits(
		diff.StartLineNumber + uint(max(len(diff.LinesToAdd) + len(diff.LinesToRemove))) - 1 + uint(len(diff.LinesAfter)),
	)

	leftPadding := strings.Repeat(" ", largestAmountOfDigits + 1)

	// Window top
	bb.WriteString(leftPadding)

	bb.WriteString(l.config.SymbolColor)
	bb.WriteString(l.config.Diff.WindowTop)
	bb.WriteString(l.config.ResetAnsi)

	bb.WriteString("\n")

	// Lines out of focus before
	for i, line := range diff.LinesBefore {
		l.renderLineNumber(
			bb,
			diff.StartLineNumber - uint(len(diff.LinesBefore)) + uint(i),
			largestAmountOfDigits,
			l.config.Diff.OutOfFocusLineNumberColor,
		)

		bb.WriteString(l.config.SymbolColor)
		bb.WriteString(l.config.Diff.LineCountSeparator)
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines to remove
	for i, line := range diff.LinesToRemove {
		l.renderLineNumber(
			bb,
			diff.StartLineNumber + uint(i),
			largestAmountOfDigits,
			l.config.Diff.InFocusLineNumberColor,
		)

		bb.WriteString(l.config.Diff.RemoveLineColor)
		bb.WriteString(l.config.Diff.RemoveLinePrefix)
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines to add
	for i, line := range diff.LinesToAdd {
		l.renderLineNumber(
			bb,
			diff.StartLineNumber + uint(i),
			largestAmountOfDigits,
			l.config.Diff.InFocusLineNumberColor,
		)

		bb.WriteString(l.config.Diff.AddLineColor)
		bb.WriteString(l.config.Diff.AddLinePrefix)
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines out of focus after
	for i, line := range diff.LinesAfter {
		l.renderLineNumber(
			bb,
			diff.StartLineNumber + uint(len(diff.LinesToAdd)) + uint(i),
			largestAmountOfDigits,
			l.config.Diff.OutOfFocusLineNumberColor,
		)

		bb.WriteString(l.config.SymbolColor)
		bb.WriteString(l.config.Diff.LineCountSeparator)
		bb.WriteString(l.config.ResetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Window bottom
	bb.WriteString(leftPadding)

	bb.WriteString(l.config.SymbolColor)
	bb.WriteString(l.config.Diff.WindowBottom)
	bb.WriteString(l.config.ResetAnsi)

	bb.WriteString("\n")
}

func countDigits(number uint) int {
	if number == 0 {
		return 1
	}

	count := 0

	for number > 0 {
		number /= 10
		count++
	}

	return count
}

func (l *LogRenderer) renderLineNumber(
	bb *bytes.Buffer,
	lineNumber uint,
	largestAmountOfDigits int,
	color string,
) {
	lineNumberAsStr := strconv.FormatUint(uint64(lineNumber), 10)

	bb.WriteString(strings.Repeat(" ", largestAmountOfDigits-len(lineNumberAsStr)))
	bb.WriteString(color)
	bb.WriteString(lineNumberAsStr)
	bb.WriteString(l.config.ResetAnsi)
	bb.WriteString(" ")
}
