package logrendering

import (
	"bytes"
	"strconv"
	"strings"
)

func countDigits(number uint) int {
	if number == 0 {
		return 1
	}

	count := 0

	for number > 0 {
		number /= 10
		count++
	}

	return count
}

func renderLineNumber(bb *bytes.Buffer, lineNumber uint, largestAmountOfDigits int, color string) {
	lineNumberAsStr := strconv.FormatUint(uint64(lineNumber), 10)

	bb.WriteString(strings.Repeat(" ", largestAmountOfDigits-len(lineNumberAsStr)))

	bb.WriteString(color)
	bb.WriteString(lineNumberAsStr)
	bb.WriteString(resetAnsi)

	bb.WriteString(" ")
}
