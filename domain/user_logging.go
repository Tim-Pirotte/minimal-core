package domain

type Severity uint8

const (
	Verbose Severity = iota
	Debug
	Info
	Warning
	SevereWarning
	Error
	Critical
)

type Message struct {
	Severity Severity
	Category string
	Message  string
}

type Annotation struct {
	Span     Span
	Message  string
	Severity Severity
}

type Line struct {
	Content     string
	Annotations []Annotation
}

type CodeContext struct {
	Source          string
	StartLineNumber uint
	LinesBefore     []string
	LinesInFocus    []Line
	LinesAfter      []string
}

type Diff struct {
	StartLineNumber uint
	LinesBefore     []string
	LinesToRemove   []string
	LinesToAdd      []string
	LinesAfter      []string
}

type Hint struct {
	Text              string
	MoreInfoReference string
}
