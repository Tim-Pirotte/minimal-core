package logrendering

import (
	"fmt"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

const (
	resetAnsi = "\033[0m"
)

func (l *LogRenderer) getSeverityColor(s usermessaging.Severity) string {
	var sAsStr string

	switch s {
	case usermessaging.Verbose:
		sAsStr = l.getStrOrDefault("severity", "verbose_color", "")
	case usermessaging.Debug:
		sAsStr = l.getStrOrDefault("severity", "debug_color", "")
	case usermessaging.Info:
		sAsStr = l.getStrOrDefault("severity", "info_color", "")
	case usermessaging.Warning:
		sAsStr = l.getStrOrDefault("severity", "warning_color", "")
	case usermessaging.SevereWarning:
		sAsStr = l.getStrOrDefault("severity", "severe_warning_color", "")
	case usermessaging.Error:
		sAsStr = l.getStrOrDefault("severity", "error_color", "")
	case usermessaging.Critical:
		sAsStr = l.getStrOrDefault("severity", "critical_color", "")
	default:
		panic(fmt.Sprintf("missing color for the enum Severity: %d", s))
	}

	return sAsStr
}

func (l *LogRenderer) getStrOrDefault(section, field, fallBack string) string {
	v, ok := l.configLoader.Get("log_renderer", section, field)

	if !ok {
		return fallBack
	}

	str, ok := v.(string)

	if !ok {
		return fallBack
	}

	return str
}

func (l *LogRenderer) getIntOrDefault(section, field string, fallBack int) int {
	v, ok := l.configLoader.Get("log_renderer", section, field)

	if !ok {
		return fallBack
	}

	number, ok := v.(int)

	if !ok {
		return fallBack
	}

	return number
}

func (l *LogRenderer) getBoolOrDefault(section, field string, fallBack bool) bool {
	v, ok := l.configLoader.Get("log_renderer", section, field)

	if !ok {
		return fallBack
	}

	b, ok := v.(bool)

	if !ok {
		return fallBack
	}

	return b
}
