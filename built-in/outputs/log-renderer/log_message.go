package logrendering

import (
	"bytes"
	messaging "minimal/minimal-core/built-in/messaging"
)

func (l *LogRenderer) outputMessage(bb *bytes.Buffer, m messaging.Message) {
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
