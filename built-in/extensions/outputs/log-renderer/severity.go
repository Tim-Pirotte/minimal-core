package formattedtext

import (
	"fmt"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

func severityToString(s usermessaging.Severity) string {
	var sAsStr string

	switch s {
	case usermessaging.Verbose:
		sAsStr = "VERBOSE"
	case usermessaging.Debug:
		sAsStr = "DEBUG"
	case usermessaging.Info:
		sAsStr = "INFO"
	case usermessaging.Warning:
		sAsStr = "WARNING"
	case usermessaging.SevereWarning:
		sAsStr = "SEVERE WARNING"
	case usermessaging.Error:
		sAsStr = "ERROR"
	case usermessaging.Critical:
		sAsStr = "CRITICAL"
	default:
		panic(fmt.Sprintf("missing string representation for the enum Severity: %d", s))
	}

	return sAsStr
}
