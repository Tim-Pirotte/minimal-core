package logrendering

import (
	"fmt"
	usermessaging "minimal/minimal-core/built-in/user-messaging"
)

const (
	resetAnsi = "\033[0m"
)

type Config struct {
	General GeneralConfig `toml:"general"`
	Severity SeverityConfig `toml:"severity"`
	Message MessageConfig `toml:"message"`
	Context ContextConfig `toml:"context"`
	Diff DiffConfig `toml:"diff"`
	InfoReference InfoReferenceConfig `toml:"info_reference"`
}

type GeneralConfig struct {
	SymbolColor     string `toml:"symbol_color"`
}

type SeverityConfig struct {
	VerboseColor       string `toml:"verbose_color"`
	DebugColor         string `toml:"debug_color"`
	InfoColor          string `toml:"info_color"`
	WarningColor       string `toml:"warning_color"`
	SevereWarningColor string `toml:"severe_warning_color"`
	ErrorColor         string `toml:"error_color"`
	CriticalColor      string `toml:"critical_color"`
}

type MessageConfig struct {
	TextHasSeverityColor bool `toml:"text_has_severity_color"`
}

type ContextConfig struct {
	SourceContextPadding int `toml:"source_context_padding"`
	TextHasSeverityColor bool `toml:"text_has_severity_color"`
	SourceStart string `toml:"source_start"`
	SourceEnd string `toml:"source_end"`
	SourceLineSep string `toml:"source_line_sep"`
	ContextWindowTopSymbol string `toml:"context_window_top_symbol"`
	ContextWindowBottomSymbol string `toml:"context_window_bottom_symbol"`
	LineCountSep string `toml:"line_count_sep"` 
	SourceColor string `toml:"source_color"`
	StartLineColor string `toml:"start_line_color"`
	AnnotationSymbol string `toml:"annotation_symbol"`
	AnnotationLineStartSymbol string `toml:"annotation_line_start_symbol"`
	AnnotationCommentStartSymbol string `toml:"annotation_comment_start_symbol"`
	OutOfFocusLineNumberColor string `toml:"out_of_focus_line_number_color"`
	InFocusLineNumberColor string `toml:"in_focus_line_number_color"`
}

type DiffConfig struct {
	DiffWindowTopSymbol string `toml:"diff_window_top_symbol"`
	DiffWindowBottomSymbol string `toml:"diff_window_bottom_symbol"`
	LineCountSep string `toml:"line_count_sep"`
	OutOfFocusLineNumberColor string `toml:"out_of_focus_line_number_color"`
	InFocusLineNumberColor string `toml:"in_focus_line_number_color"`
	RemoveLineSymbol string `toml:"remove_line_symbol"`
	AddLineSymbol string `toml:"add_line_symbol"`
	RemoveLineColor string `toml:"remove_line_color"`
	AddLineColor string `toml:"add_line_color"`
}

type InfoReferenceConfig struct {
	HintPrefixColor string `toml:"hint_prefix_color"`
	HintColor string `toml:"hint_color"`
	MoreInfoColor string `toml:"more_info_color"`
}

func (c SeverityConfig) color(s usermessaging.Severity) string {
	var sAsStr string

	switch s {
	case usermessaging.Verbose:
		sAsStr = c.VerboseColor
	case usermessaging.Debug:
		sAsStr = c.DebugColor
	case usermessaging.Info:
		sAsStr = c.InfoColor
	case usermessaging.Warning:
		sAsStr = c.WarningColor
	case usermessaging.SevereWarning:
		sAsStr = c.SevereWarningColor
	case usermessaging.Error:
		sAsStr = c.ErrorColor
	case usermessaging.Critical:
		sAsStr = c.CriticalColor
	default:
		panic(fmt.Sprintf("missing color for the enum Severity: %d", s))
	}

	return sAsStr
}
