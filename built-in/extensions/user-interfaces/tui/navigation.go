package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Navigation struct {
	panels    []tview.Primitive
}

func NewNavigation(app *tview.Application) *Navigation {
	navigation := &Navigation{make([]tview.Primitive, 0)}

	navigation.addPanelNavigation(app)
	
	return navigation
}

func (d *Navigation) AddPanel(panel tview.Primitive) {
	d.panels = append(d.panels, panel)
}

func (d *Navigation) addPanelNavigation(app *tview.Application) {
    focusIndex := 0

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			focusIndex = (focusIndex + 1) % len(d.panels)
            app.SetFocus(d.panels[focusIndex])

            return nil
		}

        return event
    })
}
