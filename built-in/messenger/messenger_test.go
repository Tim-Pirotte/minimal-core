package messenger

import (
    "strings"
    "testing"
)

func TestEmpty(t *testing.T) {
    m := New()
    m.Send(Message{})
    m.Close()
}

type orderOutput struct {
    t *testing.T
    current int
}

func (o *orderOutput) Receive(message Message) {
    if o.current != len(message.Message) {
        o.t.Fatalf("Expected message length %d but got length %d", o.current, len(message.Message))
    }

    o.current++
}

func TestOrder(t *testing.T) {
    m := New()
    m.AddOutput(&orderOutput{t, 0})

    for i := range 100 {
        m.Send(Message{Message: strings.Repeat("a", i)})
    }

    m.Close()
}

func TestSendAfterClose(t *testing.T) {
    defer func() {
        if recover() == nil {
            t.Fail()
        }
    }()

    m := New()
    m.Close()
    m.Send(Message{})
    t.Fail()
}
