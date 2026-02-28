package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {

}

func (t *TUI) StartTui(args []string) bool {
	app := tview.NewApplication()

	dashBoard := setDashBoard(app)
	app.SetRoot(dashBoard, true)

	if err := app.Run(); err != nil {
		return false
	}

	return true
}

func setDashBoard(app *tview.Application) *tview.Grid {
	actions := tview.NewList().
		AddItem("Run", "Run the project", 'r', nil).
		AddItem("Build", "Compile the project", 'b', nil).
		AddItem("Test", "Run tests", 't', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	actions.SetBorder(true).SetTitle("Actions")

	output := tview.NewTextView().
			SetBorder(true).
			SetTitle("Output")

	shellOutput := tview.NewTextView().SetChangedFunc(func() { app.Draw() })
	shellInput := tview.NewInputField().SetLabel("> ")
	shellInput.SetDoneFunc(
		func(key tcell.Key) {
			if key != tcell.KeyEnter {
				return
			}

			commandText := shellInput.GetText()
			shellInput.SetText("")

			if commandText == "" {
				return
			}

			fmt.Fprintf(shellOutput, "> %s\n", commandText)

			parts := strings.Fields(commandText)

			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Stdout = shellOutput
        	cmd.Stderr = shellOutput

			go func() {
				err := cmd.Run()

				if err != nil {
					fmt.Fprint(shellOutput, err)
				}

				shellOutput.ScrollToEnd()
			}()
		},
	)

	shell := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(shellOutput, 0, 1, false).
		AddItem(shellInput, 1, 0, true)
	
	shell.SetBorder(true).SetTitle("Shell")

	dashBoard := tview.NewGrid().
		SetRows(0).
		SetColumns(20, 0, 50)

	dashBoard.AddItem(actions, 0, 0, 1, 1, 0, 0, true).
		AddItem(output, 0, 1, 1, 1, 0, 0, false).
		AddItem(shell, 0, 2, 1, 1, 0, 0, false)

	dashBoard.SetBorder(true).SetTitle("Minimal")

	setPanelNavigation(app, []tview.Primitive{actions, output, shell})

	return dashBoard
}

func setPanelNavigation(app *tview.Application, panels []tview.Primitive) {
    focusIndex := 0

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyRight, tcell.KeyTab:
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
}
