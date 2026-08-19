package testoutput

import (
	"minimal/minimal-core/built-in/messenger"
	"reflect"
	"testing"
)

func Test(t *testing.T) {
    m := messenger.New()
    to := New()

    m.AddOutput(to)

    testMessage := messenger.Message{Message: "Test message"}

    m.Send(testMessage)
    m.Close()

    if len(to.messages) != 1 {
        t.Fatalf("Expected to.messages to have 1 message but it has %d messages", len(to.messages))
    }

    if !reflect.DeepEqual(to.messages[0], testMessage) {
        t.Error("Expected:", testMessage, "Got:", to.messages[0])
    }

    to.CheckMessages(t, []messenger.Message{testMessage})
}