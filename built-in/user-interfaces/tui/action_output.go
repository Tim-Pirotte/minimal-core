package tui

import (
	"github.com/rivo/tview"
)

type Output struct {
	window *tview.Flex
}

func NewOutput() *Output {
	return &Output{tview.NewFlex()}
}

func (o *Output) Draw(primitive tview.Primitive) {
	o.window.Clear()
	o.window.AddItem(primitive, 0, 1, true)
}
