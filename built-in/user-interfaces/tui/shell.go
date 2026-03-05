package tui

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var ErrDuplicateShellCommand = errors.New("shell command with this name already exists")

type Shell struct {
	container        *tview.Flex
	input            *tview.InputField
	output           *tview.TextView
	cmdOutput        io.Writer
	commands         map[string]func(shell *Shell, args []string)
	history          []string
	historyIndex     int
	initialDirectory string
}

func NewShell(app *tview.Application, navigation *Navigation, dashboard *tview.Grid) *Shell {
	initialDirectory, _ := os.Getwd()

	input := getShellInput()
	output := getShellOutput(app)

	shell := &Shell{
		getShellUI(input, output),
		input,
		output,
		tview.ANSIWriter(output),
		make(map[string]func(shell *Shell, args []string)), 
		make([]string, 0), 
		0,
		initialDirectory,
	}

	navigation.AddPanel(shell.input)
	navigation.AddPanel(shell.output)

	shell.addToDashboard(dashboard)
	shell.runCommandOnEnter()
	shell.addBuiltInCommands()
	shell.addHistoryNavigation()
	shell.updatePrompt()

	return shell
}

func (s *Shell) AddCommand(name string, run func(shell *Shell, args []string)) error {
	if _, ok := s.commands[name]; ok {
		return ErrDuplicateShellCommand
	}
	
	s.commands[name] = run
	
	return nil
}

func getShellInput() *tview.InputField {
	input := tview.NewInputField()

	input.SetLabelColor(tcell.ColorRed)
	input.SetBorder(true)

	return input
}

func getShellOutput(app *tview.Application) *tview.TextView {
	output := tview.NewTextView()

	output.SetChangedFunc(func() { app.Draw() })
	output.SetScrollable(true)
	output.SetDynamicColors(true)

	return output
}

func getShellUI(input *tview.InputField, output *tview.TextView) *tview.Flex {
	shell := tview.NewFlex()

	shell.SetDirection(tview.FlexRow)
	shell.AddItem(output, 0, 1, false)
	shell.AddItem(input, 3, 0, false)
	shell.SetBorder(true).SetTitle("Shell")

	return shell
}

func (s *Shell) addBuiltInCommands() {
	s.AddCommand("clear", clear)
	s.AddCommand("cd", changeDirectory)
}

func clear(shell *Shell, _ []string) {
	shell.output.Clear()
}

func changeDirectory(shell *Shell, args []string) {
	if len(args) < 2 {
		err := os.Chdir(shell.initialDirectory)

		if err != nil {
			shell.output.Write([]byte(err.Error()))
			return
		}

		shell.updatePrompt()
	} else {
		path := args[1]

		if strings.HasPrefix(path, "~") {
			home, _ := os.UserHomeDir()
			path = filepath.Join(home, path[1:])
		}

		err := os.Chdir(path)

		if err != nil {
			shell.output.Write([]byte(err.Error()))
			return
		}

		shell.updatePrompt()
	}
}

func (s *Shell) addToDashboard(dashboard *tview.Grid) {
	dashboard.AddItem(s.container, 0, 2, 1, 1, 0, 0, false)
}

func (s *Shell) addHistoryNavigation() {
	s.input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
        case tcell.KeyUp:
            if s.historyIndex > 0 {
				s.historyIndex--
				s.input.SetText(s.history[s.historyIndex])
			}

            return nil
        case tcell.KeyDown:
			if s.historyIndex < len(s.history) {
				s.historyIndex++
			}

			if s.historyIndex < len(s.history) {
				s.input.SetText(s.history[s.historyIndex])
			} else {
				s.input.SetText("")
			}
            
			return nil
        }

		return event
	})
}

func (s *Shell) runCommandOnEnter() {
	s.input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}

		commandText := s.input.GetText()
		s.input.SetText("")

		if commandText == "" {
			return
		}

		s.history = append(s.history, commandText)
		s.historyIndex = len(s.history)

		fmt.Fprintf(s.output, "\n%s%s\n", s.input.GetLabel(), commandText)

		parts := strings.Fields(commandText)

		if len(parts) > 0 {
			if run, ok := s.commands[parts[0]]; ok {
				run(s, parts)
				return
			}
		}

		s.runCommand(parts[0], parts[1:])
	})
}

func (s *Shell) runCommand(name string, args []string) {
	cmd := exec.Command(name, args...)

	cmd.Stdout = s.cmdOutput
	cmd.Stderr = s.cmdOutput

	go func() {
		err := cmd.Run()

		if err != nil {
			fmt.Fprint(s.cmdOutput, err)
		}

		s.output.ScrollToEnd()
	}()
}

func (s *Shell) updatePrompt() {
	path, err := os.Getwd()

	if err != nil {
		s.output.Write([]byte(err.Error()))
		return
	}

	s.input.SetLabel(path + "|>")
}
