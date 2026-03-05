package tui

import (
	"github.com/rivo/tview"
)

type SimpleAction struct {
	Name        string
	Description string
	ShortCut    rune
	FnToRun     func(output *Output)
}

func (s *SimpleAction) GetName() string {
	return s.Name
}

func (s *SimpleAction) GetDescription() string {
	return s.Description
}

func (s *SimpleAction) GetShortCut() rune {
	return s.ShortCut
}

func (s *SimpleAction) Run(output *Output) {
	if s.FnToRun != nil {
		s.FnToRun(output)
	} else {
		displayNotImplemented(output)
	}
}

func displayNotImplemented(output *Output) {
	output.Draw(tview.NewTextView().SetText("This action doesn't have an implementation"))
}
