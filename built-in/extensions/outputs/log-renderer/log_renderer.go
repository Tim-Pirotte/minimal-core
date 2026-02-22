package logrendering

import (
	"bytes"
	"fmt"
	"io"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

type LogRenderer struct {
	writer io.Writer
	config Config
}

// Use this function for a new logger but **don't forget** to close the logger.
// This is necessary to ensure that all the messages are displayed
func NewLogger(writer io.Writer, config Config) *LogRenderer {
	return &LogRenderer{writer, config}
}

type bytesBuffer struct {
	bytesBuffer *bytes.Buffer
}

func (b *bytesBuffer) Handle() {}

func (l *LogRenderer) CreateHandle() usermessaging.Handle {
	return &bytesBuffer{bytes.NewBuffer(make([]byte, 0))}
}

func (l *LogRenderer) Finish(h usermessaging.Handle) {
	b, _ := h.(*bytesBuffer)
	bb := b.bytesBuffer

	if bb == nil {
		panic("an uninitialized bytes buffer was passed to the log displayer")
	}

	l.writer.Write([]byte(resetAnsi))
	_, err := l.writer.Write(bb.Bytes())

	if err != nil {
		// TODO error message
		fmt.Println("Logging error:", err)
	}
}
