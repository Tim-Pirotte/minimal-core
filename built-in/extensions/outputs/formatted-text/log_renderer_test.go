package formattedtext

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

func TestDisplayRenderLogAtomic(t *testing.T) {
	configFile, _ := os.ReadFile("./config.toml")
	config, _ := LoadConfig(configFile)

	finalBuffer := &bytes.Buffer{}
	logger := NewLogger(finalBuffer, config)
	
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
			logger.FlushOutput(buf)
		}
	}()

	go func() {
		defer wg.Done()
		for range iterations {
			buf := bytes.NewBufferString(strB)
			logger.FlushOutput(buf)
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

		if len(line) != stringLength {
			t.Errorf("Atomicity failure at line %d: expected length %d, got %d", lineCount, stringLength, len(line))
		}
	}

	expectedTotalLines := iterations * 2
	if lineCount != expectedTotalLines {
		t.Errorf("Expected %d total lines, but got %d", expectedTotalLines, lineCount)
	}
}

func BenchmarkLogger(b *testing.B) {
	configFile, _ := os.ReadFile("./config.toml")
	config, _ := LoadConfig(configFile)

	logger := NewLogger(io.Discard, config)
	defer logger.Close()

	ctx := CodeContext{"main.go", 10, []string{}, []Line{{Content: "fmt.Println(\"hello\")"}}, []string{}}

	m := Message{Info, "BENCHMARK", "Testing logger performance"}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
			bb := bytes.Buffer{}
            logger.OutputMessage(&bb, m)
			logger.OutputContext(&bb, ctx)
			logger.FlushOutput(&bb)
        }
    })
}
