package userlogging

import "minimal/minimal-core/domain"

type Logger struct {
	outputs []Output
}

type Output interface {
	OutputMessage(domain.Message)
	OutputContext(domain.CodeContext)
	OutputDiff(domain.Diff)
	OutputHint(domain.Hint)
}

func New() *Logger {
	return &Logger{
		outputs: make([]Output, 0),
	}
}

func (l *Logger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
}

func (l *Logger) LogMessage(message domain.Message) {
	for _, o := range l.outputs {
		o.OutputMessage(message)
	}
}

func (l *Logger) LogContext(context domain.CodeContext) {
	for _, o := range l.outputs {
		o.OutputContext(context)
	}
}

func (l *Logger) LogDiff(diff domain.Diff) {
	for _, o := range l.outputs {
		o.OutputDiff(diff)
	}
}

func (l *Logger) LogHint(hint domain.Hint) {
	for _, o := range l.outputs {
		o.OutputHint(hint)
	}
}
