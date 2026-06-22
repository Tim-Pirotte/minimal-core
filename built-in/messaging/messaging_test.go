package messaging

import (
	"fmt"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	"sync"
	"testing"
)

// TODO Add proper tests
func TestMessenger(t *testing.T) {
	sourceGen := logging.GetTestLogSource(io.Discard)
	m := NewMessenger(sourceGen)
	mock := &TestOutput{}
	m.AddOutput(mock)

	m.Send([]MessagePart{
		&Message{Critical, "Hello, World!"},
	})

	m.Close()

	expectedLength := 1

	if len(mock.messages) != expectedLength {
		t.Fatal("Expected length of messages to be", expectedLength, "but got", len(mock.messages))
	}

	expectedInnerLength := 1

	if len(mock.messages[0]) != expectedInnerLength {
		t.Fatal(
			"Expected length of the message to be",
			expectedInnerLength,
			"but got",
			len(mock.messages[0]),
		)
	}

	message, ok := mock.messages[0][0].(*Message)

	if !ok {
		t.Fatal(
			"Expected the first message to be of type *Message but got",
			fmt.Sprintf("%T\n", mock.messages[0][0]),
		)
	}

	expectedMessage := "Hello, World!"

	if message.Message != expectedMessage {
		t.Fatal("Expected message to be", expectedMessage, "but got", message.Message)
	}
}

func TestMessengerRaceConditions(t *testing.T) {
	sourceGen := logging.GetTestLogSource(io.Discard)
	m := NewMessenger(sourceGen)
	mock := &TestOutput{}
	m.AddOutput(mock)

	const goroutines = 10
	const messagesPerRoutine = 100
	var wg sync.WaitGroup

	for i := range goroutines {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			messages := make([]MessagePart, messagesPerRoutine)

			for i := range messagesPerRoutine {
				messages[i] = &Message{Critical, "Hello, World!"}
			}

			m.Send(messages)
		}(i)
	}

	wg.Wait()
	m.Close()

	expectedTotal := goroutines * messagesPerRoutine

	actualLength := 0

	for _, messageParts := range mock.messages {
		actualLength += len(messageParts)
	}

	if actualLength != expectedTotal {
		t.Errorf("Data loss detected: expected %d messages, got %d", expectedTotal, len(mock.messages))
	}
}
