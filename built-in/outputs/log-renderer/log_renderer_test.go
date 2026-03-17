package logrendering

import (
	"bufio"
	"bytes"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
	"strings"
	"sync"
	"testing"
)

func TestDisplayRenderLogAtomic(t *testing.T) {
	config := Config{}

	finalBuffer := &bytes.Buffer{}
	logger := NewLogRenderer(&logging.SourceGenerator{}, finalBuffer, config)
	
	iterations := 1000
	stringLength := 100
	
	strA := strings.Repeat("A", stringLength) + "\n"
	strB := strings.Repeat("B", stringLength) + "\n"

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for range iterations {
			buf := bytes.NewBufferString(strA)
			handle := &bytesBuffer{buf}
			logger.Finish(handle)
		}
	}()

	go func() {
		defer wg.Done()
		for range iterations {
			buf := bytes.NewBufferString(strB)
			handle := &bytesBuffer{buf}
			logger.Finish(handle)
		}
	}()

	wg.Wait()

	scanner := bufio.NewScanner(finalBuffer)
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		hasA := strings.Contains(line, "A")
		hasB := strings.Contains(line, "B")

		if hasA && hasB {
			t.Errorf("Atomicity failure at line %d: mixed content detected: %s", lineCount, line)
		}

		if len(line) != stringLength + len(resetAnsi) {
			t.Errorf("Atomicity failure at line %d: expected length %d, got %d", lineCount, stringLength + len(resetAnsi), len(line))
		}
	}

	expectedTotalLines := iterations * 2
	if lineCount != expectedTotalLines {
		t.Errorf("Expected %d total lines, but got %d", expectedTotalLines, lineCount)
	}
}

func BenchmarkLogger(b *testing.B) {
	config := Config{}

	logger := NewLogRenderer(logging.GetTestLogSource(io.Discard), io.Discard, config)

	ctx := usermessaging.CodeContext{
		Source: "main.go", 
		StartLineNumber: 10, 
		LinesBefore: []string{}, 
		LinesInFocus: []usermessaging.Line{{Content: "fmt.Println(\"hello\")"}}, 
		LinesAfter: []string{},
	}

	m := usermessaging.Message{
		Severity: usermessaging.Info, 
		Category: "BENCHMARK", 
		Message: "Testing logger performance",
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
			bb := bytes.Buffer{}
			handle := &bytesBuffer{&bb}
            logger.OutputMessage(handle, m)
			logger.OutputContext(handle, ctx)
			logger.Finish(handle)
        }
    })
}
