package formattedtext

import "bytes"

// A hint provided to the user with an optional reference for more info
// The reference can be left out by passing an empty string
func (l *LogRenderer) RenderLogHint(bb *bytes.Buffer, hint Hint) {
	if bb == nil {
        panic("an uninitialized bytes buffer was passed to the hint renderer")
    }

	bb.WriteString(l.config.InfoReference.HintPrefixColor)
	bb.WriteString("Hint")
	bb.WriteString(resetAnsi)

	bb.WriteString(l.config.General.SymbolColor)
	bb.WriteString(": ")
	bb.WriteString(resetAnsi)

	bb.WriteString(l.config.InfoReference.HintColor)
	bb.WriteString(hint.Text)
	bb.WriteString(resetAnsi)
	bb.WriteString("\n")

	if hint.MoreInfoReference != "" {
		bb.WriteString(l.config.InfoReference.MoreInfoColor)
		bb.WriteString("More info about this")
		bb.WriteString(resetAnsi)

		bb.WriteString(l.config.General.SymbolColor)
		bb.WriteString(": ")
		bb.WriteString(resetAnsi)

		bb.WriteString(hint.MoreInfoReference)
		bb.WriteString("\n")
	}
}

type Hint struct {
	Text string
	MoreInfoReference string
}
