package tui

import (
	configloader "minimal/minimal-core/built-in/config-loader"

	"github.com/rivo/tview"
)

type TUI struct {
	app            *tview.Application
	actions        *Actions
	shell          *Shell
}

func NewTUI() *TUI {
	app := tview.NewApplication()
	dashboard := tview.NewGrid().SetRows(0).SetColumns(30, 0, 40)
	navigation := NewNavigation(app)

	tui := &TUI{
		app, 
		NewActions(app, navigation, dashboard), 
		NewShell(app, navigation, dashboard),
	}

	app.SetRoot(dashboard, true)

	return tui
}

func (t *TUI) StartTUI(configLoader configloader.ConfigLoader, args []string) bool {
	if err := t.app.Run(); err != nil {
		return false
	}

	return true
}
