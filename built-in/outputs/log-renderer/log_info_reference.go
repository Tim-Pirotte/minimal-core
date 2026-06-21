package logrendering

import (
	"bytes"
	messaging "minimal/minimal-core/built-in/messaging"
)

func (l *LogRenderer) outputHint(bb *bytes.Buffer, hint messaging.Hint) {
	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the hint renderer")
    }

	if hint.Text != "" {
		bb.WriteString(l.getStrOrDefault("info_reference", "hint_prefix_color", ""))
		bb.WriteString("Hint")
		bb.WriteString(resetAnsi)

		bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
		bb.WriteString(": ")
		bb.WriteString(resetAnsi)

		bb.WriteString(l.getStrOrDefault("info_reference", "hint_color", ""))
		bb.WriteString(hint.Text)
		bb.WriteString(resetAnsi)
		bb.WriteString("\n")
	}

	if hint.MoreInfoReference != "" {
		bb.WriteString(l.getStrOrDefault("info_reference", "more_info_color", ""))
		bb.WriteString("More info about this")
		bb.WriteString(resetAnsi)

		bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
		bb.WriteString(": ")
		bb.WriteString(resetAnsi)

		bb.WriteString(hint.MoreInfoReference)
		bb.WriteString("\n")
	}
}
