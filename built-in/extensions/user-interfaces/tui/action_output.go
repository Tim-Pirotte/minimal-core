package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Output struct {
	*tview.Box
	RenderFunc func(screen tcell.Screen, x, y, w, h int)
}

func NewOutput() *Output {
	return &Output{tview.NewBox(), nil}
}

func (o *Output) Draw(screen tcell.Screen) {
	o.Box.DrawForSubclass(screen, o)
	x, y, w, h := o.GetInnerRect()

	if o.RenderFunc == nil {
		return
	}

	o.RenderFunc(screen, x, y, w, h)
}
