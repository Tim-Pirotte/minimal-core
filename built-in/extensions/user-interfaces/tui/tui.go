package tui

import "github.com/rivo/tview"

type TUI struct {

}

func (t *TUI) StartTui(args []string) bool {
	app := tview.NewApplication()

	newGridItem := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	actions := tview.NewList().
		AddItem("Run", "Run the project", 'r', nil).
		AddItem("Build", "Compile the project", 'b', nil).
		AddItem("Test", "Run tests", 't', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	actions.SetBorder(true).SetTitle("Actions")

	currentApp := newGridItem("Current App")
	shell := newGridItem("Shell")

	grid := tview.NewGrid().
		SetRows(0).
		SetColumns(20, 0, 20).
		SetBorders(true)

	grid.AddItem(actions, 0, 0, 1, 1, 0, 0, false).
		AddItem(currentApp, 0, 1, 1, 1, 0, 0, false).
		AddItem(shell, 0, 2, 1, 1, 0, 0, false)

	grid.SetBorder(true).SetTitle("Minimal")

	app.SetRoot(grid, true)
	app.SetFocus(actions)

	if err := app.Run(); err != nil {
		return false
	}

	return true
}
