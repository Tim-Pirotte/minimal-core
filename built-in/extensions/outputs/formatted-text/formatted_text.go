package formattedtext

import (
	"bytes"
	"io"
	usermessaging "minimal/minimal-core/built-in/user-messaging"

	logRenderer "github.com/Tim-Pirotte/log-renderer/log-renderer"
)

type FormattedTextMessages struct {
	logRenderer *logRenderer.LogRenderer
	currentMessageBuffer *bytes.Buffer
}

func NewFormattedTextMessages() *FormattedTextMessages {
	return &FormattedTextMessages{logRenderer.NewLogger(io.Discard, logRenderer.Config{}), bytes.NewBuffer(make([]byte, 0))}
}

func (f *FormattedTextMessages) OutputMessage(message usermessaging.Message) {
	f.logRenderer.RenderLogMessage(
		f.currentMessageBuffer, 
		logRenderer.Message{
			Severity: mapSeverity(message.Severity), 
			Category: message.Category, 
			Message: message.Message,
		},
	)
}

func mapSeverity(s usermessaging.Severity) logRenderer.Severity {
	var targetSeverity logRenderer.Severity
	
	switch s {
	case usermessaging.Debug:
		targetSeverity = logRenderer.Debug
	case usermessaging.Verbose:
		targetSeverity = logRenderer.Verbose
	case usermessaging.Info:
		targetSeverity = logRenderer.Info
	case usermessaging.Warning:
		targetSeverity = logRenderer.Warning
	case usermessaging.SevereWarning:
		targetSeverity = logRenderer.SevereWarning
	case usermessaging.Error:
		targetSeverity = logRenderer.Error
	case usermessaging.Critical:
		targetSeverity = logRenderer.Critical
	default:
		panic("unreachable")
	}

	return targetSeverity
}

func (f *FormattedTextMessages) OutputContext(codeContext usermessaging.CodeContext) {
	f.logRenderer.RenderLogContext(
		f.currentMessageBuffer, 
		logRenderer.CodeContext{
			Source: codeContext.Source, 
			StartLineNumber: codeContext.StartLineNumber,
			LinesBefore: codeContext.LinesBefore,
			LinesInFocus: mapLines(codeContext.LinesInFocus),
			LinesAfter: codeContext.LinesAfter,
		},
	)
}

func mapLines(l []usermessaging.Line) []logRenderer.Line {
	result := make([]logRenderer.Line, len(l))

	for i := range l {
		annotations := make([]logRenderer.Annotation, len(l[i].Annotations))

		for j, a := range l[i].Annotations {
			annotations[j] = logRenderer.Annotation{
				Start: int(a.Span.Start), 
				Length: int(a.Span.Length), 
				Message: a.Message, 
				Severity: mapSeverity(a.Severity),
			}
		}

		result[i] = logRenderer.Line{Content: l[i].Content, Annotations: annotations}
	}

	return result
}

func (f *FormattedTextMessages) OutputDiff(diff usermessaging.Diff) {

}

func (f *FormattedTextMessages) OutputHint(hint usermessaging.Hint) {

}