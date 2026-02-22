package logrendering

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"strings"
)

// A part of the source code with proposed edits
func (l *LogRenderer) OutputDiff(h usermessaging.Handle, diff usermessaging.Diff) {
	b, _ := h.(*bytesBuffer)
	bb := b.bytesBuffer

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

	bb.WriteString(l.config.General.SymbolColor)
	bb.WriteString(l.config.Diff.DiffWindowTopSymbol)
	bb.WriteString(resetAnsi)

	bb.WriteString("\n")

	// Lines out of focus before
	for i, line := range diff.LinesBefore {
		renderLineNumber(bb, diff.StartLineNumber - uint(len(diff.LinesBefore)) + uint(i), largestAmountOfDigits, l.config.Diff.OutOfFocusLineNumberColor)
		
		bb.WriteString(l.config.General.SymbolColor)
		bb.WriteString(l.config.Diff.LineCountSep)
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}
	
	// Lines to remove
	for i, line := range diff.LinesToRemove {
		renderLineNumber(bb, diff.StartLineNumber + uint(i), largestAmountOfDigits, l.config.Diff.InFocusLineNumberColor)
		
		bb.WriteString(l.config.Diff.RemoveLineColor)
		bb.WriteString(l.config.Diff.RemoveLineSymbol)
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines to add
	for i, line := range diff.LinesToAdd {
		renderLineNumber(bb, diff.StartLineNumber + uint(i), largestAmountOfDigits, l.config.Diff.InFocusLineNumberColor)

		bb.WriteString(l.config.Diff.AddLineColor)
		bb.WriteString(l.config.Diff.AddLineSymbol)
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines out of focus after
	for i, line := range diff.LinesAfter {
		renderLineNumber(bb, diff.StartLineNumber + uint(len(diff.LinesToAdd)) + uint(i), largestAmountOfDigits, l.config.Diff.OutOfFocusLineNumberColor)
		
		bb.WriteString(l.config.General.SymbolColor)
		bb.WriteString(l.config.Diff.LineCountSep)
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Window bottom
	bb.WriteString(leftPadding)

	bb.WriteString(l.config.General.SymbolColor)
	bb.WriteString(l.config.Diff.DiffWindowBottomSymbol)
	bb.WriteString(resetAnsi)

	bb.WriteString("\n")
}

type Diff struct {
	StartLineNumber uint
	LinesBefore   []string
	LinesToRemove []string
	LinesToAdd    []string
	LinesAfter    []string
}
