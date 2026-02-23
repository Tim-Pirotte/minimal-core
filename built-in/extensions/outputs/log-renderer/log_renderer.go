package logrendering

import (
	"bytes"
	"io"
	logging "minimal/minimal-core/built-in/internal-logging"
	usermessaging "minimal/minimal-core/built-in/user-messaging"

	"github.com/rs/zerolog"
)

type LogRenderer struct {
	logger zerolog.Logger
	writer io.Writer
	config Config
}

func NewLogger(sourceGen *logging.SourceGenerator, writer io.Writer, config Config) *LogRenderer {
	logger, _ := sourceGen.GetLogger("logRendering")

	return &LogRenderer{logger, writer, config}
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

	l.writer.Write([]byte(resetAnsi))
	_, err := l.writer.Write(bb.Bytes())

	if err != nil {
		l.logger.Error().Err(err).Msg("unsuccessfull write")
	}
}
