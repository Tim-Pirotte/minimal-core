package internallogging

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger 

var activeBuffer *ringBuffer

func New(maxBytes uint) Logger {
	activeBuffer = newRingBuffer(maxBytes)

    return zerolog.New(activeBuffer).With().Timestamp().Logger()
}

func WriteTo(w io.Writer) (int64, error) {
    if activeBuffer == nil {
        panic("the ring buffer for the logger has not been initialized")
    }

    return activeBuffer.writeTo(w)
}
