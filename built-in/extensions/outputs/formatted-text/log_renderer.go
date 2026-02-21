package formattedtext

import (
	"bytes"
	"fmt"
	"io"
)

type LogRenderer struct {
	writer io.Writer
	queue  chan *bytes.Buffer
	done   chan bool
	config Config
}

// Use this function for a new logger but **don't forget** to close the logger.
// This is necessary to ensure that all the messages are displayed
func NewLogger(writer io.Writer, config Config) *LogRenderer {
	l := LogRenderer{writer, make(chan *bytes.Buffer, 10), make(chan bool), config}

	go l.startConsumer()

	return &l
}

func (l *LogRenderer) startConsumer() {
	for log := range l.queue {
		if log == nil {
        	panic("an uninitialized bytes buffer was passed to the log displayer")
    	}

		log.WriteString(resetAnsi)

		_, err := l.writer.Write(log.Bytes())
	
		if err != nil {
			fmt.Println("Loggin error:", err)
		}
	}

	l.done <- true
}

func (l *LogRenderer) DisplayRenderedLog(bb *bytes.Buffer) {
	l.queue <- bb
}

// Call this function after you are done with the logger
func (l *LogRenderer) Close() {
    close(l.queue)
    <-l.done
}
