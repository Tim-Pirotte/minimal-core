package logrendering

import (
	"bytes"
	"io"
	configloader "minimal/minimal-core/built-in/config-loader"
	logging "minimal/minimal-core/built-in/internal-logging"
	"minimal/minimal-core/built-in/messaging"
)

type LogRenderer struct {
	logger logging.Logger
	writer io.Writer
	configLoader configloader.ConfigLoader
}

func NewLogRenderer(
	sourceGen *logging.SourceGenerator,
	writer io.Writer,
	configLoader configloader.ConfigLoader,
) *LogRenderer {
	logger, _ := sourceGen.GetLogger("logRendering")

	return &LogRenderer{logger, writer, configLoader}
}

func (l *LogRenderer) Output(messageParts []messaging.MessageType) {
	bb := bytes.NewBuffer(make([]byte, 0))

	for _, part := range messageParts {
		switch p := part.(type) {
		case *messaging.Message:
			l.OutputMessage(bb, *p)
		case *messaging.CodeContext:
			l.OutputContext(bb, *p)
		case *messaging.Hint:
			l.OutputHint(bb, *p)
		case *messaging.Diff:
			l.OutputDiff(bb, *p)
		default:
			// TODO change to proper error
			panic("unsuported message part")
		}
	}

	l.writer.Write([]byte(resetAnsi))
	_, err := l.writer.Write(bb.Bytes())

	if err != nil {
		l.logger.Error().Err(err).Msg("unsuccessfull write")
	}
}
