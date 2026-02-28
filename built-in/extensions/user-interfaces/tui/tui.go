package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type TUI struct {
	app            *tview.Application
	actions        *tview.List
	commandHistory []string
	historyIndex   int
}

type Action struct {
	Name        string
	Description string
	Shortcut    rune
	Run         func()
}

func NewTUI() *TUI {
	tui := &TUI{tview.NewApplication(), tview.NewList(), make([]string, 0), 0}
	
	tui.AddAction(Action{"Run", "Run the project", 'r', nil})
	tui.AddAction(Action{"Build", "Compile the project", 'b', nil})
	tui.AddAction(Action{"Test", "Run tests", 't', nil})
	tui.AddAction(Action{"Quit", "Press to exit", 'q', func() {
		tui.app.Stop()
	}})

	return tui
}

func (t *TUI) AddAction(action Action) {
	t.actions.AddItem(action.Name, action.Description, action.Shortcut, action.Run)
}

func (t *TUI) StartTUI(args []string) bool {
	dashBoard := t.setDashBoard()
	t.app.SetRoot(dashBoard, true)

	if err := t.app.Run(); err != nil {
		return false
	}

	return true
}

func (t *TUI) setDashBoard() *tview.Grid {
	t.actions.SetBorder(true).SetTitle("Actions")

	output := tview.NewTextView().
			SetBorder(true).
			SetTitle("Output")

	shellOutput := tview.NewTextView().SetChangedFunc(func() { t.app.Draw() })
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

			t.commandHistory = append(t.commandHistory, commandText)
			t.historyIndex = len(t.commandHistory)

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

	shell.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
        case tcell.KeyUp:
            if t.historyIndex > 0 {
				t.historyIndex--
				shellInput.SetText(t.commandHistory[t.historyIndex])
			}

            return nil
        case tcell.KeyDown:
			if t.historyIndex < len(t.commandHistory) {
				t.historyIndex++
			}

			if t.historyIndex < len(t.commandHistory) {
				shellInput.SetText(t.commandHistory[t.historyIndex])
			} else {
				shellInput.SetText("")
			}
            
			return nil
        }

		return event
	})

	dashBoard := tview.NewGrid().
		SetRows(0).
		SetColumns(20, 0, 50)

	dashBoard.AddItem(t.actions, 0, 0, 1, 1, 0, 0, true).
		AddItem(output, 0, 1, 1, 1, 0, 0, false).
		AddItem(shell, 0, 2, 1, 1, 0, 0, false)

	dashBoard.SetBorder(true).SetTitle("Minimal")

	t.setPanelNavigation([]tview.Primitive{t.actions, output, shell})

	return dashBoard
}

func (t *TUI) setPanelNavigation(panels []tview.Primitive) {
    focusIndex := 0

	t.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        switch event.Key() {
        case tcell.KeyRight, tcell.KeyTab:
            focusIndex = (focusIndex + 1) % len(panels)
            t.app.SetFocus(panels[focusIndex])

            return nil
        case tcell.KeyLeft:
            focusIndex = (focusIndex - 1 + len(panels)) % len(panels)
            t.app.SetFocus(panels[focusIndex])
            
			return nil
        }

        return event
    })
}
