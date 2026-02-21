package formattedtext

import "bytes"

// A simple message
func (l *LogRenderer) RenderLogMessage(bb *bytes.Buffer, m Message) {
	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the message renderer")
    }

	bb.WriteString(l.config.Severity.color(m.Severity))
	bb.WriteString(m.Severity.string())
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

type Message struct {
	Severity Severity
	Category string
	Message  string
}
