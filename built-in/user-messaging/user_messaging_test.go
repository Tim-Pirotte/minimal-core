package usermessaging

import (
	"sync"
	"testing"
)

type MockHandle struct{
	messages []Message
}

func (m *MockHandle) Handle() {}

type MockOutput struct {
	sync.Mutex
	messages []Message
	finished bool
}

func (m *MockOutput) CreateHandle() Handle { 
	return &MockHandle{} 
}

func (m *MockOutput) Finish(h Handle) { 
	m.Lock()
	defer m.Unlock()
	
	mockHandle, _ := h.(*MockHandle)
	m.messages = append(m.messages, mockHandle.messages...)
	m.finished = true
}

func (m *MockOutput) OutputMessage(h Handle, msg Message) {
	mockHandle, _ := h.(*MockHandle)
	mockHandle.messages = append(mockHandle.messages, msg)
}

func (m *MockOutput) OutputContext(h Handle, ctx CodeContext) {}
func (m *MockOutput) OutputDiff(h Handle, d Diff)             {}
func (m *MockOutput) OutputHint(h Handle, hi Hint)            {}

func TestMessenger(t *testing.T) {
	m := NewMessenger()
	mock := &MockOutput{}
	m.AddOutput(mock)

	tx := m.CreateLogTransaction()
	m.LogMessage(tx, Message{Critical, "Category", "Hello, World!"})
	m.CommitLogTransaction(tx)
	m.Close()

	expectedLength := 1

	if len(mock.messages) != expectedLength {
		t.Fatal("Expected length of messages to be", expectedLength, "but got", len(mock.messages))
	}

	expectedCategory := "Category"

	if mock.messages[0].Category != expectedCategory {
		t.Fatal("Expected category to be", expectedCategory, "but got", mock.messages[0].Category)
	}

	expectedMessages := "Hello, World!"

	if mock.messages[0].Message != expectedMessages {
		t.Fatal("Expected message to be", expectedMessages, "but got", mock.messages[0].Message)
	}

	if !mock.finished {
		t.Error("Expected transaction to be finished")
	}
}

func TestMessengerRaceConditions(t *testing.T) {
	m := NewMessenger()
	mock := &MockOutput{}
	m.AddOutput(mock)

	const goroutines = 10
	const messagesPerRoutine = 100
	var wg sync.WaitGroup

	for i := range goroutines {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			tx := m.CreateLogTransaction()
			
			for range messagesPerRoutine {
				m.LogMessage(tx, Message{Critical, "Category", "Hello, World!"})
			}
			
			m.CommitLogTransaction(tx)
		}(i)
	}

	wg.Wait()
	m.Close()

	expectedTotal := goroutines * messagesPerRoutine

	if len(mock.messages) != expectedTotal {
		t.Errorf("Data loss detected: expected %d messages, got %d", expectedTotal, len(mock.messages))
	}
}
