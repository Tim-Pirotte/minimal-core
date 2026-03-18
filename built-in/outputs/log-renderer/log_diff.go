package logrendering

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"strings"
)

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

	bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
	bb.WriteString(l.getStrOrDefault("diff", "window_top_symbol", ""))
	bb.WriteString(resetAnsi)

	bb.WriteString("\n")

	// Lines out of focus before
	for i, line := range diff.LinesBefore {
		renderLineNumber(
			bb, 
			diff.StartLineNumber - uint(len(diff.LinesBefore)) + uint(i), 
			largestAmountOfDigits, 
			l.getStrOrDefault("diff", "out_of_focus_line_number_color", ""),
		)
		
		bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
		bb.WriteString(l.getStrOrDefault("diff", "line_count_separator", "|"))
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}
	
	// Lines to remove
	for i, line := range diff.LinesToRemove {
		renderLineNumber(
			bb, 
			diff.StartLineNumber + uint(i), 
			largestAmountOfDigits, 
			l.getStrOrDefault("diff", "in_focus_line_number_color", ""),
		)
		
		bb.WriteString(l.getStrOrDefault("diff", "remove_line_color", ""))
		bb.WriteString(l.getStrOrDefault("diff", "remove_line_symbol", "-"))
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines to add
	for i, line := range diff.LinesToAdd {
		renderLineNumber(
			bb, 
			diff.StartLineNumber + uint(i), 
			largestAmountOfDigits, 
			l.getStrOrDefault("diff", "in_focus_line_number_color", ""),
		)

		bb.WriteString(l.getStrOrDefault("diff", "add_line_color", ""))
		bb.WriteString(l.getStrOrDefault("diff", "add_line_symbol", "+"))
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Lines out of focus after
	for i, line := range diff.LinesAfter {
		renderLineNumber(
			bb, 
			diff.StartLineNumber + uint(len(diff.LinesToAdd)) + uint(i), 
			largestAmountOfDigits, 
			l.getStrOrDefault("diff", "out_of_focus_line_number_color", ""),
		)
		
		bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
		bb.WriteString(l.getStrOrDefault("diff", "line_count_separator", "|"))
		bb.WriteString(resetAnsi)

		bb.WriteString(" ")
		bb.WriteString(line)
		bb.WriteString("\n")
	}

	// Window bottom
	bb.WriteString(leftPadding)

	bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
	bb.WriteString(l.getStrOrDefault("diff", "window_bottom_symbol", ""))
	bb.WriteString(resetAnsi)

	bb.WriteString("\n")
}
