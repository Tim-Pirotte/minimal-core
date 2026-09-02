package testoutput

import (
	"fmt"
	"minimal/minimal-core/built-in/messenger"
	logrendering "minimal/minimal-core/built-in/outputs/log-renderer"
	"minimal/minimal-core/built-in/substring"
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

        t.Fail()
    } else {
        for i := range expected {
            for j := range expected[i].Context {
                compareSubStrings(
                    m,
                    t,
                    expected[i].Context[j].Content,
                    to.messages[i].Context[j].Content,
                    expected[i].Message,
                    i + 1,
                )
            }

            for j := range expected[i].AdditionalContext {
                compareSubStrings(
                    m,
                    t,
                    expected[i].AdditionalContext[j].Content,
                    to.messages[i].AdditionalContext[j].Content,
                    expected[i].Message,
                    i + 1,
                )
            }

            for j := range expected[i].Suggestions {
                for k := range expected[i].Suggestions[j].Replacements {
                    compareSubStrings(
                        m,
                        t,
                        expected[i].Suggestions[j].Replacements[k].From.Content,
                        to.messages[i].Suggestions[j].Replacements[k].From.Content,
                        expected[i].Message,
                        i + 1,
                    )
                }
            }
        }
    }

    m.Close()
}

func compareSubStrings(m *messenger.Messenger, t *testing.T, expected, actual, message string, number int) {
    if !substring.IsSubString(expected, actual) {
        m.Send(
            messenger.Message{
                Message: "Different address than expected",
                Severity: messenger.Error,
                Notes: []string{
                    fmt.Sprintf("Message number: %d", number),
                    fmt.Sprintf("Message text: %s", message),
                    fmt.Sprintf("Content: %s", expected),
                },
            },
        )

        t.Fail()
    }
}

func (t *TestOutput) Receive(message messenger.Message) {
    t.messages = append(t.messages, message)
}
