package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Actions struct {
	actions *tview.List
	output  *Output
}

type Action interface {
	GetName() string
	GetDescription() string
	GetShortCut() rune
	Run(output *Output)
}

func NewActions(app *tview.Application, navigation *Navigation, dashboard *tview.Grid) *Actions {
	actions := &Actions{tview.NewList(), NewOutput()}

	navigation.AddPanel(actions.actions)
	navigation.AddPanel(actions.output.window)

	actions.addToDashboard(dashboard)
	actions.addStyling()
	actions.addBuiltInActions(app)
	actions.addSelectionHighlight()

	return actions
}

func (a *Actions) AddAction(action Action) {
	a.actions.AddItem(
		action.GetName(),
		action.GetDescription(), 
		action.GetShortCut(), 
		func() { action.Run(a.output) },
	)
}

func (a *Actions) addToDashboard(dashboard *tview.Grid) {
	dashboard.AddItem(a.actions, 0, 0, 1, 1, 0, 0, true)
	dashboard.AddItem(a.output.window, 0, 1, 1, 1, 0, 0, false)
}

func (a *Actions) addStyling() {
	a.actions.SetBorder(true).SetTitle("Actions")
	a.output.window.SetBorder(true).SetTitle("Output")
}

func (a *Actions) addBuiltInActions(app *tview.Application) {
	a.AddAction(&SimpleAction{"Run", "Run the project", 'r', nil})
	a.AddAction(&SimpleAction{"Build", "Compile the project", 'b', nil})
	a.AddAction(&SimpleAction{"Test", "Run tests", 't', nil})
	a.AddAction(&SimpleAction{"Quit", "Press to exit", 'q', func(output *Output) {
		app.Stop()
	}})
}

func (a *Actions) addSelectionHighlight() {
	a.actions.SetFocusFunc(func() {
		a.actions.SetSelectedBackgroundColor(tcell.ColorBlue)
		a.actions.SetSelectedTextColor(tcell.ColorWhite)
	})

	a.actions.SetBlurFunc(func() {
		a.actions.SetSelectedBackgroundColor(tcell.ColorDefault)
    	a.actions.SetSelectedTextColor(tcell.ColorWhite)
	})
}
