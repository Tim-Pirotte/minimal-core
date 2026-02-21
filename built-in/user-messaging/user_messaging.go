package usermessaging

type Messenger struct {
	outputs []Output
}

type Transaction struct {
	handles map[Output]Handle
}

type Output interface {
	CreateHandle() Handle
	Finish(Handle)
	OutputMessage(Handle, Message)
	OutputContext(Handle, CodeContext)
	OutputDiff(Handle, Diff)
	OutputHint(Handle, Hint)
}

type Handle interface {
	IsHandle()
}

func NewMessenger() *Messenger {
	return &Messenger{
		outputs: make([]Output, 0),
	}
}

func (l *Messenger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
}

func (l *Messenger) CreateLogTransaction() *Transaction {
	transaction := Transaction{make(map[Output]Handle)}

	for _, o := range l.outputs {
		transaction.handles[o] = o.CreateHandle()
	}

	return &transaction
}

func (l *Messenger) CommitLogTransaction(transaction *Transaction) {
	for _, o := range l.outputs {
		o.Finish(transaction.handles[o])
	}
}

func (l *Messenger) LogMessage(transaction *Transaction, message Message) {
	for _, o := range l.outputs {
		o.OutputMessage(transaction.handles[o], message)
	}
}

func (l *Messenger) LogContext(transaction *Transaction, context CodeContext) {
	for _, o := range l.outputs {
		o.OutputContext(transaction.handles[o], context)
	}
}

func (l *Messenger) LogDiff(transaction *Transaction, diff Diff) {
	for _, o := range l.outputs {
		o.OutputDiff(transaction.handles[o], diff)
	}
}

func (l *Messenger) LogHint(transaction *Transaction, hint Hint) {
	for _, o := range l.outputs {
		o.OutputHint(transaction.handles[o], hint)
	}
}
