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
	queue   chan []MessageType
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
	Output([]MessageType)
}

type MessageType interface {
	MessageType()
}

type Message struct {
	Severity Severity
	Category string
	Message  string
}

func (*Message) MessageType() {}

type CodeContext struct {
	Source          string
	StartLineNumber uint
	LinesBefore     []string
	LinesInFocus    []Line
	LinesAfter      []string
}

func (*CodeContext) MessageType() {}

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

func (*Diff) MessageType() {}

type Hint struct {
	Text              string
	MoreInfoReference string
}

func (*Hint) MessageType() {}

func NewMessenger(sourceGen *logging.SourceGenerator) *Messenger {
	logger, _ := sourceGen.GetLogger("messenger")

	m := &Messenger{
		logger:  logger,
		outputs: make([]Output, 0),
		queue:   make(chan []MessageType, bufferSize),
		done:    make(chan bool),
	}

	go m.worker()

	return m
}

func (m *Messenger) worker() {
	for messages := range m.queue {
		for _, output := range m.outputs {
			output.Output(messages)
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

func (m *Messenger) Output(messages []MessageType) {
	m.queue<-messages
}
