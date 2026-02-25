package usermessaging

const bufferSize = 10

type Messenger struct {
	outputs []Output
	queue   chan func()
	done    chan bool
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
	Handle()
}

func NewMessenger() *Messenger {
	m := &Messenger{
		outputs: make([]Output, 0),
		queue:   make(chan func(), bufferSize),
		done:    make(chan bool),
	}

	go m.worker()

	return m
}

func (m *Messenger) worker() {
	for job := range m.queue {
		job()
	}

	close(m.done)
}

func (l *Messenger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
}

func (m *Messenger) Close() {
	close(m.queue)
	<-m.done
}

func (l *Messenger) CreateLogTransaction() *Transaction {
	transaction := Transaction{make(map[Output]Handle)}

	for _, o := range l.outputs {
		transaction.handles[o] = o.CreateHandle()
	}

	return &transaction
}

func (m *Messenger) CommitLogTransaction(transaction *Transaction) {
	m.runPerOutput(transaction, func(o Output, h Handle) { o.Finish(h) })
}

func (m *Messenger) LogMessage(transaction *Transaction, message Message) {
	m.runPerOutput(transaction, func(o Output, h Handle) { o.OutputMessage(h, message) })
}

func (m *Messenger) LogContext(transaction *Transaction, context CodeContext) {
	m.runPerOutput(transaction, func(o Output, h Handle) { o.OutputContext(h, context) })
}

func (m *Messenger) LogDiff(transaction *Transaction, diff Diff) {
	m.runPerOutput(transaction, func(o Output, h Handle) { o.OutputDiff(h, diff) })
}

func (m *Messenger) LogHint(transaction *Transaction, hint Hint) {
	m.runPerOutput(transaction, func(o Output, h Handle) { o.OutputHint(h, hint) })
}

func (m *Messenger) runPerOutput(transaction *Transaction, f func(Output, Handle)) {
	m.queue <- func() {
		for _, o := range m.outputs {
			if handle, ok := transaction.handles[o]; ok {
				f(o, handle)
			}
		}
	}
}
