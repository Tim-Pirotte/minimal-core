package logrendering

import (
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

func (l *LogRenderer) OutputMessage(h usermessaging.Handle, m usermessaging.Message) {
	b, _ := h.(*bytesBuffer)
	bb := b.bytesBuffer

	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the message renderer")
    }

	bb.WriteString(l.getSeverityColor(m.Severity))
	bb.WriteString(severityToString(m.Severity))
	bb.WriteString(resetAnsi)
	bb.WriteString(" ")
	bb.WriteString(m.Category)
	bb.WriteString(l.getStrOrDefault("general", "symbol_color", ""))
	bb.WriteString(": ")
	bb.WriteString(resetAnsi)

	if l.getBoolOrDefault("general", "text_has_severity_color", false) {
		bb.WriteString(l.getSeverityColor(m.Severity))
	}
	
	bb.WriteString(m.Message)

	if l.getBoolOrDefault("general", "text_has_severity_color", false) {
		bb.WriteString(resetAnsi)
	}
	
	bb.WriteString("\n")
}
