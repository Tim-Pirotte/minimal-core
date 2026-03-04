package tui

import (
	"github.com/gdamore/tcell/v2"
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
	output.RenderFunc = func(screen tcell.Screen, x, y, w, h int) {
		tview.Print(
			screen, 
			"This action doesn't have an implementation", 
			x,
			y,
			w, 
			tview.AlignCenter, 
			tcell.ColorWhite,
		)
	}
}
