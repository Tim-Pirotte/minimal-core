package messaging

import (
	"fmt"
	logging "minimal/minimal-core/built-in/internal-logging"
	"reflect"
	"testing"
)

const bufferSize = 16

type Messenger struct {
	logger  logging.Logger
	outputs []Output
	queue   chan Message
	done    chan bool
}

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

type Message struct {
	Reference         string
	Message           string
	Severity          Severity
	Context           []Span
	AdditionalContext []Span
	Suggestions       []Suggestion
	Notes             []string
}

type Span struct {
	Content string
	Note    string
}

type Suggestion struct {
	Suggestion string
	Replacements []Replacement
}

type Replacement struct {
	From Span
	To   Span
}

type Output interface {
	Receive(Message)
}

type TestOutput struct {
	messages []Message
}

func (m *TestOutput) Receive(message Message) {
	m.messages = append(m.messages, message)
}

func NewMessenger(sourceGen *logging.SourceGenerator) *Messenger {
	logger, _ := sourceGen.GetLogger("Messenger")

	m := &Messenger{
		logger:  logger,
		outputs: make([]Output, 0),
		queue:   make(chan Message, bufferSize),
		done:    make(chan bool),
	}

	go m.worker()

	return m
}

func (m *Messenger) worker() {
	for message := range m.queue {
		for _, output := range m.outputs {
			output.Receive(message)
		}
	}

	close(m.done)
}

func (l *Messenger) AddOutput(outputChannel Output) {
	l.outputs = append(l.outputs, outputChannel)
	l.logger.Debug().Str("output", fmt.Sprintf("%v", reflect.TypeOf(outputChannel))).Msg("output registered")
}

func (m *Messenger) Close() {
	close(m.queue)
	<-m.done
}

func (m *Messenger) Send(message Message) {
	m.queue<-message
}

func (to *TestOutput) CheckMessages(t *testing.T, expected []Message) {
	// TODO make this more sophisticated
	if !reflect.DeepEqual(to.messages, expected) {
		t.Error("Expected:\n", expected, "\nGot:\n", to.messages)
	}
}
