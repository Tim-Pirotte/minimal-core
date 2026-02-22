package logrendering

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

// A simple message
func (l *LogRenderer) OutputMessage(h usermessaging.Handle, m usermessaging.Message) {
	b, _ := h.(*bytesBuffer)
	bb := b.bytesBuffer

	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the message renderer")
    }

	bb.WriteString(l.config.Severity.color(m.Severity))
	bb.WriteString(severityToString(m.Severity))
	bb.WriteString(resetAnsi)
	bb.WriteString(" ")
	bb.WriteString(m.Category)
	bb.WriteString(l.config.General.SymbolColor)
	bb.WriteString(": ")
	bb.WriteString(resetAnsi)

	if l.config.Message.TextHasSeverityColor {
		bb.WriteString(l.config.Severity.color(m.Severity))
	}
	
	bb.WriteString(m.Message)

	if l.config.Message.TextHasSeverityColor {
		bb.WriteString(resetAnsi)
	}
	
	bb.WriteString("\n")
}
