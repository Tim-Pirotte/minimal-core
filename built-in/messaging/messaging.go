package messaging

import (
	"reflect"
	"testing"
)

// TODO allow stacktraces

const bufferSize = 16

type Messenger struct {
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

func NewMessenger() *Messenger {
    m := &Messenger{
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
}

func (m *Messenger) Close() {
    close(m.queue)
    <-m.done
}

func (m *Messenger) Send(message Message) {
    m.queue<-message
}

func (to *TestOutput) CheckMessages(t *testing.T, expected []Message) {
    if to.messages == nil {
        to.messages = []Message{}
    }

    if expected == nil {
        expected = []Message{}
    }

    // TODO make this more sophisticated
    if !reflect.DeepEqual(to.messages, expected) {
        t.Error("\nExpected:\n", expected, "\nGot:\n", to.messages)
    }
}
