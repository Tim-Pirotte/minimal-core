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

	_, err := rb.Write([]byte(firstInput))

	if err != nil {
		t.Error("Write to ringbuffer failed")
	}

	var firstOutput bytes.Buffer

	_, err = rb.WriteTo(&firstOutput)

	if err != nil {
		t.Error("Ringbuffer WriteTo failed")
	}

	firstExpected := "Hello, "

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	_, err = rb.Write([]byte(secondInput))

	if err != nil {
		t.Error("Write to ringbuffer failed")
	}

	var secondOutput bytes.Buffer

	_, err = rb.WriteTo(&secondOutput)

	if err != nil {
		t.Error("Ringbuffer WriteTo failed")
	}

	secondExpected := "lo, world!"

	if secondOutput.String() != secondExpected {
		t.Error("Expected", secondExpected, "but got", secondOutput.String())
	}
}

func TestEmptyBuffer(t *testing.T) {
	resetInit()

	rb := NewRingBuffer(0)

	firstInput := "Hello, "

	_, err := rb.Write([]byte(firstInput))

	if err != nil {
		t.Error("Write to ringbuffer failed")
	}

	var firstOutput bytes.Buffer

	_, err = rb.WriteTo(&firstOutput)

	if err != nil {
		t.Error("ringbuffer WriteTo failed")
	}

	firstExpected := ""

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	_, err = rb.Write([]byte(secondInput))

	if err != nil {
		t.Error("Write to ringbuffer failed")
	}

	var secondOutput bytes.Buffer

	_, err = rb.WriteTo(&secondOutput)

	if err != nil {
		t.Error("ringbuffer WriteTo failed")
	}

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

	firstLogger, firstSource := rootSource.GetLogger("FirstLevel")
	secondLogger, _ := firstSource.GetLogger("SecondLevel")

	firstLogger.Info().Msg("Hello, world!")
	secondLogger.Info().Msg("Hello, world!")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	firstExpected := "\"source\":[\"FirstLevel\",\"SecondLevel\"]"

	if !strings.Contains(actual, firstExpected) {
		t.Error("Expected", firstExpected, "in", actual)
	}

	secondExpected := "\"source\":[\"FirstLevel\"]"

	if !strings.Contains(actual, secondExpected) {
		t.Error("Expected", secondExpected, "in", actual)
	}
}

func TestDuplicateSources(t *testing.T) {
	resetInit()

	ringBuffer := NewRingBuffer(500)
	rootSource := Init(ringBuffer)

	firstLogger, _ := rootSource.GetLogger("Duplicate")
	secondLogger, _ := rootSource.GetLogger("Duplicate")

	firstLogger.Info().Msg("Hello, world!")
	secondLogger.Info().Msg("Hello, world!")

	var buf bytes.Buffer
	_, err := ringBuffer.WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	firstExpected := "\"source\":[\"Duplicate\"]"

	if !strings.Contains(actual, firstExpected) {
		t.Error("Expected", firstExpected, "in", actual)
	}

	secondExpected := "\"source\":[\"Duplicate#1\"]"

	if !strings.Contains(actual, secondExpected) {
		t.Error("Expected", secondExpected, "in", actual)
	}
}
