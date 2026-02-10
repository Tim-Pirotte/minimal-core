package logging

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestRingBuffer(t *testing.T) {
	rb := newRingBuffer(10)

	firstInput := "Hello, "

	rb.Write([]byte(firstInput))

	var firstOutput bytes.Buffer
	
	rb.writeTo(&firstOutput)

	firstExpected := "Hello, "

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	rb.Write([]byte(secondInput))

	var secondOutput bytes.Buffer
	
	rb.writeTo(&secondOutput)

	secondExpected := "lo, world!"

	if secondOutput.String() != secondExpected {
		t.Error("Expected", secondExpected, "but got", secondOutput.String())
	}
}

func TestEmptyBuffer(t *testing.T) {
	rb := newRingBuffer(0)

	firstInput := "Hello, "

	rb.Write([]byte(firstInput))

	var firstOutput bytes.Buffer
	
	rb.writeTo(&firstOutput)

	firstExpected := ""

	if firstOutput.String() != firstExpected {
		t.Error("Expected", firstExpected, "but got", firstOutput.String())
	}

	secondInput := "world!"

	rb.Write([]byte(secondInput))

	var secondOutput bytes.Buffer
	
	rb.writeTo(&secondOutput)

	secondExpected := ""

	if secondOutput.String() != secondExpected {
		t.Error("Expected", secondExpected, "but got", secondOutput.String())
	}
}

func TestLogger(t *testing.T) {
	rootSource := Init(500)

	logger, _ := rootSource.GetLogger("")

	logger.Info().Msg("Hello, world!")
	logger.Error().Msg("Error message")

	var buf bytes.Buffer
	_, err := WriteTo(&buf)

	if err != nil {
		t.Fatalf("Failed to write to buffer: %v", err)
	}

	actual := buf.String()

	fmt.Println(actual)

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
