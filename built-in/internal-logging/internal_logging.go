package internallogging

import (
	"io"

	"github.com/rs/zerolog"
)

type Logger = zerolog.Logger 

func New(w io.Writer) Logger {
    return zerolog.New(w).With().Timestamp().Logger()
}
