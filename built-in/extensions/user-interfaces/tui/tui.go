package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {

}

func (t *TUI) StartTui(args []string) bool {
	app := tview.NewApplication()

	actions := tview.NewList().
		AddItem("Run", "Run the project", 'r', nil).
		AddItem("Build", "Compile the project", 'b', nil).
		AddItem("Test", "Run tests", 't', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	actions.SetBorder(true).SetTitle("Actions")

	currentApp := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("Current App")

	shell := tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText("Shell")

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

	panels := []tview.Primitive{actions, currentApp, shell}
    focusIndex := 0

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyRight:
            focusIndex = (focusIndex + 1) % len(panels)
            app.SetFocus(panels[focusIndex])

            return nil
        case tcell.KeyLeft:
            focusIndex = (focusIndex - 1 + len(panels)) % len(panels)
            app.SetFocus(panels[focusIndex])
            
			return nil
        }

        return event
    })

	if err := app.Run(); err != nil {
		return false
	}

	return true
}
