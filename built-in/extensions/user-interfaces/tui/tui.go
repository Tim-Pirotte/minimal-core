package tui

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ErrDuplicateShellCommand = errors.New("shell command with this name already exists")

type TUI struct {
	app            *tview.Application
	actions        *tview.List
	shellOutput    *tview.TextView
	shellCommands  map[string]func(args []string)
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
	tui := &TUI{
		tview.NewApplication(), 
		tview.NewList(), 
		tview.NewTextView(),
		make(map[string]func(args []string)), 
		make([]string, 0), 
		0,
	}
	
	tui.AddAction(Action{"Run", "Run the project", 'r', nil})
	tui.AddAction(Action{"Build", "Compile the project", 'b', nil})
	tui.AddAction(Action{"Test", "Run tests", 't', nil})
	tui.AddAction(Action{"Quit", "Press to exit", 'q', func() {
		tui.app.Stop()
	}})

	tui.AddShellCommand("clear", func(args []string) {
		tui.shellOutput.Clear()
	})

	return tui
}

func (t *TUI) AddAction(action Action) {
	t.actions.AddItem(action.Name, action.Description, action.Shortcut, action.Run)
}

func (t *TUI) AddShellCommand(name string, run func(args []string)) error {
	if _, ok := t.shellCommands[name]; ok {
		return ErrDuplicateShellCommand
	}
	
	t.shellCommands[name] = run
	
	return nil
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

	t.actions.SetFocusFunc(func() {
		t.actions.SetSelectedBackgroundColor(tcell.ColorBlue)
		t.actions.SetSelectedTextColor(tcell.ColorWhite)
	})

	t.actions.SetBlurFunc(func() {
		t.actions.SetSelectedBackgroundColor(tcell.ColorDefault)
    	t.actions.SetSelectedTextColor(tcell.ColorWhite)
	})

	output := tview.NewTextView()
	output.SetBorder(true).SetTitle("Output")

	t.shellOutput.SetChangedFunc(func() { t.app.Draw() })
	t.shellOutput.SetScrollable(true)
	t.shellOutput.SetDynamicColors(true)

	shellInput := tview.NewInputField().SetLabel(":> ")
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

			fmt.Fprintf(t.shellOutput, "\n|> %s\n", commandText)

			parts := strings.Fields(commandText)

			if len(parts) > 0 {
				if run, ok := t.shellCommands[parts[0]]; ok {
					run(parts)
					return
				}
			}

			ansiWriter := tview.ANSIWriter(t.shellOutput)
			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Stdout = ansiWriter
        	cmd.Stderr = ansiWriter

			go func() {
				err := cmd.Run()

				if err != nil {
					fmt.Fprint(ansiWriter, err)
				}

				t.shellOutput.ScrollToEnd()
			}()
		},
	)

	shellInput.SetLabelColor(tcell.ColorWhite)

	shellInput.SetFocusFunc(func() {
		shellInput.SetLabelColor(tcell.ColorRed)
		shellInput.SetLabel("|> ")
	})

	shellInput.SetBlurFunc(func() {
		shellInput.SetLabelColor(tcell.ColorWhite)
		shellInput.SetLabel(":> ")
	})

	shell := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.shellOutput, 0, 1, false).
		AddItem(shellInput, 1, 0, false)
	
	shell.SetBorder(true).SetTitle("Shell")

	shellInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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
		SetColumns(30, 0, 40)

	dashBoard.AddItem(t.actions, 0, 0, 1, 1, 0, 0, true).
		AddItem(output, 0, 1, 1, 1, 0, 0, false).
		AddItem(shell, 0, 2, 1, 1, 0, 0, false)

	dashBoard.SetBorder(true).SetTitle("Minimal")

	t.setPanelNavigation([]tview.Primitive{t.actions, output, shellInput, t.shellOutput})

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
