package logrendering

import (
	"fmt"
	"minimal/minimal-core/built-in/messaging"
)

const (
	resetAnsi = "\033[0m"
)

func (l *LogRenderer) getSeverityColor(s messaging.Severity) string {
	var sAsStr string

	switch s {
	case messaging.Verbose:
		sAsStr = l.getStrOrDefault("severity", "verbose_color", "")
	case messaging.Debug:
		sAsStr = l.getStrOrDefault("severity", "debug_color", "")
	case messaging.Info:
		sAsStr = l.getStrOrDefault("severity", "info_color", "")
	case messaging.Warning:
		sAsStr = l.getStrOrDefault("severity", "warning_color", "")
	case messaging.SevereWarning:
		sAsStr = l.getStrOrDefault("severity", "severe_warning_color", "")
	case messaging.Error:
		sAsStr = l.getStrOrDefault("severity", "error_color", "")
	case messaging.Critical:
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
