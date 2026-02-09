package internallogging

import (
	"errors"
	"io"

	"github.com/rs/zerolog"
)

const UninitializedRingBuffer = "the ring buffer for the logger has not been initialized"

type Logger = zerolog.Logger 

var activeBuffer *ringBuffer

func New(maxBytes uint) Logger {
	activeBuffer = newRingBuffer(maxBytes)

    return zerolog.New(activeBuffer).With().Timestamp().Logger()
}

func WriteTo(w io.Writer) (int64, error) {
    if activeBuffer == nil {
        return 0, errors.New(UninitializedRingBuffer)
    }

    return activeBuffer.writeTo(w)
}
