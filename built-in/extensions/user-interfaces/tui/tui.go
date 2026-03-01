package tui

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ErrDuplicateShellCommand = errors.New("shell command with this name already exists")

type TUI struct {
	app            *tview.Application
	actions        *tview.List
	shell          shell
}

type shell struct {
	input            *tview.InputField
	output           *tview.TextView
	commands         map[string]func(args []string)
	history          []string
	historyIndex     int
	initialDirectory string
}

type Action struct {
	Name        string
	Description string
	Shortcut    rune
	Run         func()
}

func NewTUI() *TUI {
	initialDirectory, _ := os.Getwd()

	tui := &TUI{
		tview.NewApplication(), 
		tview.NewList(), 
		shell{
			tview.NewInputField(),
			tview.NewTextView(),
			make(map[string]func(args []string)), 
			make([]string, 0), 
			0,
			initialDirectory,
		},
	}
	
	tui.AddAction(Action{"Run", "Run the project", 'r', nil})
	tui.AddAction(Action{"Build", "Compile the project", 'b', nil})
	tui.AddAction(Action{"Test", "Run tests", 't', nil})
	tui.AddAction(Action{"Quit", "Press to exit", 'q', func() {
		tui.app.Stop()
	}})

	tui.AddShellCommand("clear", func(args []string) {
		tui.shell.output.Clear()
	})

	tui.AddShellCommand("cd", func(args []string) {
		if len(args) < 2 {
			err := os.Chdir(tui.shell.initialDirectory)

			if err != nil {
				tui.shell.output.Write([]byte(err.Error()))
				return
			}

			tui.shell.input.SetLabel(tui.shell.initialDirectory + " |> ")
		} else {
			path := args[1]

			if strings.HasPrefix(path, "~") {
				home, _ := os.UserHomeDir()
				path = filepath.Join(home, path[1:])
			}

			err := os.Chdir(path)

			if err != nil {
				tui.shell.output.Write([]byte(err.Error()))
				return
			}

			tui.shell.input.SetLabel(path + " |> ")
		}
	})

	return tui
}

func (t *TUI) AddAction(action Action) {
	t.actions.AddItem(action.Name, action.Description, action.Shortcut, action.Run)
}

func (t *TUI) AddShellCommand(name string, run func(args []string)) error {
	if _, ok := t.shell.commands[name]; ok {
		return ErrDuplicateShellCommand
	}
	
	t.shell.commands[name] = run
	
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

	t.shell.output.SetChangedFunc(func() { t.app.Draw() })
	t.shell.output.SetScrollable(true)
	t.shell.output.SetDynamicColors(true)

	t.shell.input.SetLabel("|  ")
	t.shell.input.SetDoneFunc(
		func(key tcell.Key) {
			if key != tcell.KeyEnter {
				return
			}

			commandText := t.shell.input.GetText()
			t.shell.input.SetText("")

			if commandText == "" {
				return
			}

			t.shell.history = append(t.shell.history, commandText)
			t.shell.historyIndex = len(t.shell.history)

			fmt.Fprintf(t.shell.output, "\n|> %s\n", commandText)

			parts := strings.Fields(commandText)

			if len(parts) > 0 {
				if run, ok := t.shell.commands[parts[0]]; ok {
					run(parts)
					return
				}
			}

			ansiWriter := tview.ANSIWriter(t.shell.output)
			cmd := exec.Command(parts[0], parts[1:]...)
			cmd.Stdout = ansiWriter
        	cmd.Stderr = ansiWriter

			go func() {
				err := cmd.Run()

				if err != nil {
					fmt.Fprint(ansiWriter, err)
				}

				t.shell.output.ScrollToEnd()
			}()
		},
	)

	t.shell.input.SetLabelColor(tcell.ColorWhite)

	t.shell.input.SetFocusFunc(func() {
		t.shell.input.SetLabelColor(tcell.ColorRed)
		t.shell.input.SetLabel("|> ")
	})

	t.shell.input.SetBlurFunc(func() {
		t.shell.input.SetLabelColor(tcell.ColorWhite)
		t.shell.input.SetLabel("|  ")
	})

	shell := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(t.shell.output, 0, 1, false).
		AddItem(t.shell.input, 1, 0, false)
	
	shell.SetBorder(true).SetTitle("Shell")

	t.shell.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
        case tcell.KeyUp:
            if t.shell.historyIndex > 0 {
				t.shell.historyIndex--
				t.shell.input.SetText(t.shell.history[t.shell.historyIndex])
			}

            return nil
        case tcell.KeyDown:
			if t.shell.historyIndex < len(t.shell.history) {
				t.shell.historyIndex++
			}

			if t.shell.historyIndex < len(t.shell.history) {
				t.shell.input.SetText(t.shell.history[t.shell.historyIndex])
			} else {
				t.shell.input.SetText("")
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

	t.setPanelNavigation([]tview.Primitive{t.actions, output, t.shell.input, t.shell.output})

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
