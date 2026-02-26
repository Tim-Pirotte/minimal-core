package tui

import "github.com/rivo/tview"

type TUI struct {

}

func (t *TUI) StartTui(args []string) bool {
	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")

	if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
		return false
	}

	return true
}
