package logging

import (
	"bytes"
	"strings"
	"testing"
)

func resetInit() {
	initCalled = false
}

func TestRingBuffer(t *testing.T) {
	resetInit()
	
	rb := NewRingBuffer(10)

	firstInput := "Hello, "

	rb.Write([]byte(firstInput))

	var firstOutput bytes.Buffer
	
	rb.WriteTo(&firstOutput)

	firstExpected := "Hello, "

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	rb.Write([]byte(secondInput))

	var secondOutput bytes.Buffer
	
	rb.WriteTo(&secondOutput)

	secondExpected := "lo, world!"

	if secondOutput.String() != secondExpected {
		t.Error("Expected", secondExpected, "but got", secondOutput.String())
	}
}

func TestEmptyBuffer(t *testing.T) {
	resetInit()

	rb := NewRingBuffer(0)

	firstInput := "Hello, "

	rb.Write([]byte(firstInput))

	var firstOutput bytes.Buffer
	
	rb.WriteTo(&firstOutput)

	firstExpected := ""

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	rb.Write([]byte(secondInput))

	var secondOutput bytes.Buffer
	
	rb.WriteTo(&secondOutput)

	secondExpected := ""

	if secondOutput.String() != secondExpected {
		t.Error("Expected", secondExpected, "but got", secondOutput.String())
	}
}

func TestLogger(t *testing.T) {
	resetInit()

	ringBuffer := NewRingBuffer(500)
	rootSource := Init(ringBuffer)

	logger, _ := rootSource.GetLogger("")

	logger.Info().Msg("Hello, world!")
	logger.Error().Msg("Error message")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	firstExpected := "\"message\":\"Hello, world!\""

	if !strings.Contains(actual, firstExpected) {
		t.Error("Expected", firstExpected, "in", actual)
	}

	secondExpected := "\"message\":\"Error message\""

	if !strings.Contains(actual, secondExpected) {
		t.Error("Expected", secondExpected, "in", actual)
	}

	thirdExpected := "\"level\":\"info\""
	
	if !strings.Contains(actual, thirdExpected) {
		t.Error("Expected", thirdExpected, "in", actual)
	}

	fourthExpected := "\"level\":\"error\""

	if !strings.Contains(actual, fourthExpected) {
		t.Error("Expected", fourthExpected, "in", actual)
	}
}

func TestEmptySource(t *testing.T) {
	resetInit()

	ringBuffer := NewRingBuffer(500)
	rootSource := Init(ringBuffer)

	logger, _ := rootSource.GetLogger("")

	logger.Info().Msg("Hello, world!")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	expected := "\"source\":[\"unnamed\"]"

	if !strings.Contains(actual, expected) {
		t.Error("Expected", expected, "in", actual)
	}
}

func TestMultipleSources(t *testing.T) {
	resetInit()

	ringBuffer := NewRingBuffer(500)
	rootSource := Init(ringBuffer)

	firstLogger, firstSource := rootSource.GetLogger("firstLevel")
	secondLogger, _ := firstSource.GetLogger("secondLevel")

	firstLogger.Info().Msg("Hello, world!")
	secondLogger.Info().Msg("Hello, world!")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	firstExpected := "\"source\":[\"firstLevel\",\"secondLevel\"]"

	if !strings.Contains(actual, firstExpected) {
		t.Error("Expected", firstExpected, "in", actual)
	}

	secondExpected := "\"source\":[\"firstLevel\"]"

	if !strings.Contains(actual, secondExpected) {
		t.Error("Expected", secondExpected, "in", actual)
	}
}

func TestDuplicateSources(t *testing.T) {
	resetInit()

	ringBuffer := NewRingBuffer(500)
	rootSource := Init(ringBuffer)

	firstLogger, _ := rootSource.GetLogger("duplicate")
	secondLogger, _ := rootSource.GetLogger("duplicate")

	firstLogger.Info().Msg("Hello, world!")
	secondLogger.Info().Msg("Hello, world!")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	firstExpected := "\"source\":[\"duplicate\"]"

	if !strings.Contains(actual, firstExpected) {
		t.Error("Expected", firstExpected, "in", actual)
	}

	secondExpected := "\"source\":[\"duplicate#1\"]"

	if !strings.Contains(actual, secondExpected) {
		t.Error("Expected", secondExpected, "in", actual)
	}
}
