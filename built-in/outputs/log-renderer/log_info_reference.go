package logrendering

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

func (l *LogRenderer) OutputHint(handle usermessaging.Handle, hint usermessaging.Hint) {
	b, _ := handle.(*bytesBuffer)
	bb := b.bytesBuffer
	
	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the hint renderer")
    }

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
