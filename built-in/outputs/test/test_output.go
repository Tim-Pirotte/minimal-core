package testoutput

import (
	"minimal/minimal-core/built-in/messenger"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"os"
	"reflect"
	"testing"
)

type TestOutput struct {
	messages []messenger.Message
	Output   messenger.Output
}

func New() *TestOutput {
    return &TestOutput{[]messenger.Message{}, logrendering.New(os.Stdout)}
}

func (to *TestOutput) CheckMessages(t *testing.T, expected []messenger.Message) {
    if to.messages == nil {
        to.messages = []messenger.Message{}
    }

    if expected == nil {
        expected = []messenger.Message{}
    }

    m := messenger.New()
    m.AddOutput(to.Output)

    if !reflect.DeepEqual(to.messages, expected) {
        m.Send(
            messenger.Message{
                Message: "TestOutput messages are not the ones expected",
                Severity: messenger.Error,
            },
        )

        m.Send(messenger.Message{Message: "Expected the following messages", Severity: messenger.Info})

        for _, message := range expected {
            m.Send(message)
        }

        m.Send(messenger.Message{Message: "Received the following messages", Severity: messenger.Info})

        for _, message := range to.messages {
            m.Send(message)
        }

        m.Close()
        t.Fail()
    }
}

func (t *TestOutput) Receive(message messenger.Message) {
    t.messages = append(t.messages, message)
}
