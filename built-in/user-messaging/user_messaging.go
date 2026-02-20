package usermessaging

type Messenger struct {
	outputs []Output
}

type Output interface {
	OutputMessage(Message)
	OutputContext(CodeContext)
	OutputDiff(Diff)
	OutputHint(Hint)
}

func New() *Messenger {
	return &Messenger{
		outputs: make([]Output, 0),
	}
}

func (l *Messenger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
}

func (l *Messenger) LogMessage(message Message) {
	for _, o := range l.outputs {
		o.OutputMessage(message)
	}
}

func (l *Messenger) LogContext(context CodeContext) {
	for _, o := range l.outputs {
		o.OutputContext(context)
	}
}

func (l *Messenger) LogDiff(diff Diff) {
	for _, o := range l.outputs {
		o.OutputDiff(diff)
	}
}

func (l *Messenger) LogHint(hint Hint) {
	for _, o := range l.outputs {
		o.OutputHint(hint)
	}
}
