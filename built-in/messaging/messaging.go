package messaging

import (
	"fmt"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/primitives"
	"reflect"
)

const bufferSize = 16

type Messenger struct {
	logger  logging.Logger
	outputs []Output
	queue   chan []MessagePart
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

type Output interface {
	Receive([]MessagePart)
}

type TestOutput struct {
	messages [][]MessagePart
}

func (m *TestOutput) Receive(messageParts []MessagePart) {
	m.messages = append(m.messages, messageParts)
}

type MessagePart interface {
	MessagePart()
}

type Message struct {
	Severity Severity
	Message  string
}

func (*Message) MessagePart() {}

type CodeContext struct {
	Source          string
	StartLineNumber uint
	LinesBefore     []string
	LinesInFocus    []Line
	LinesAfter      []string
}

func (*CodeContext) MessagePart() {}

type Line struct {
	Content     string
	Annotations []Annotation
}

type Annotation struct {
	Range    primitives.Range
	Message  string
	Severity Severity
}

type Diff struct {
	StartLineNumber uint
	LinesBefore     []string
	LinesToRemove   []string
	LinesToAdd      []string
	LinesAfter      []string
}

func (*Diff) MessagePart() {}

type Hint struct {
	Text              string
	MoreInfoReference string
}

func (*Hint) MessagePart() {}

func NewMessenger(sourceGen *logging.SourceGenerator) *Messenger {
	logger, _ := sourceGen.GetLogger("Messenger")

	m := &Messenger{
		logger:  logger,
		outputs: make([]Output, 0),
		queue:   make(chan []MessagePart, bufferSize),
		done:    make(chan bool),
	}

	go m.worker()

	return m
}

func (m *Messenger) worker() {
	for messageParts := range m.queue {
		for _, output := range m.outputs {
			output.Receive(messageParts)
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

func (m *Messenger) Send(messageParts []MessagePart) {
	m.queue<-messageParts
}
