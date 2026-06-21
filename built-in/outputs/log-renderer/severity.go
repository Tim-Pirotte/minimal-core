package logrendering

import (
	"fmt"
	messaging "minimal/minimal-core/built-in/messaging"
)

func severityToString(s messaging.Severity) string {
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
