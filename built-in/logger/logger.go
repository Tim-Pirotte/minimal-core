package logger

type Logger struct {
	outputs []Output
}

type Output interface {
}

func New() *Logger {
	return &Logger{
		outputs: make([]Output, 0),
	}
}

func (l *Logger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
}

type Message struct {
}

func (l *Logger) LogMessage() {

}

type Context struct {
}

func (l *Logger) LogContext() {

}

type Diff struct {
}

func (l *Logger) LogDiff() {

}

type Hint struct {
}

func (l *Logger) LogHint() {

}
