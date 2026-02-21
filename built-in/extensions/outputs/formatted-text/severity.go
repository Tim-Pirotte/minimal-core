package formattedtext

import "fmt"

// The importance level of a message
type Severity uint8

const (
	Verbose Severity = iota
	Debug
	Info
	Warning
	SevereWarning
	Error
	Critical
)

func (s Severity) string() string {
	var sAsStr string

	switch s {
	case Verbose:
		sAsStr = "VERBOSE"
	case Debug:
		sAsStr = "DEBUG"
	case Info:
		sAsStr = "INFO"
	case Warning:
		sAsStr = "WARNING"
	case SevereWarning:
		sAsStr = "SEVERE WARNING"
	case Error:
		sAsStr = "ERROR"
	case Critical:
		sAsStr = "CRITICAL"
	default:
		panic(fmt.Sprintf("missing string representation for the enum Severity: %d", s))
	}

	return sAsStr
}
