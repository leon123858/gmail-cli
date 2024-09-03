package tui

import (
	"github.com/rivo/tview"
)

type App struct {
	App    *tview.Application
	Pages  *tview.Pages
	Widget []tview.Primitive
}

func NewApp() *App {
	return &App{
		App:    tview.NewApplication(),
		Pages:  tview.NewPages(),
		Widget: []tview.Primitive{},
	}
}

func (ta *App) GetPage() *tview.Pages {
	return ta.Pages
}

func (ta *App) AddPage(name string, item tview.Primitive) {
	ta.Pages.AddPage(name, item, true, false)
}

func (ta *App) ShowPage(name string) {
	ta.Pages.SwitchToPage(name)
}

func (ta *App) AddWidget(w tview.Primitive) {
	ta.Widget = append(ta.Widget, w)
}

func (ta *App) GetCurrentFocus() tview.Primitive {
	return ta.App.GetFocus()
}

func (ta *App) SetNextFocus() {
	widget := ta.GetCurrentFocus()
	ta.App.SetFocus(ta.NextWidgets(widget))
}

func (ta *App) PreviousWidgets(current tview.Primitive) tview.Primitive {
	for i, w := range ta.Widget {
		if w == current {
			if i-1 >= 0 {
				return ta.Widget[i-1]
			}
			return ta.Widget[len(ta.Widget)-1]
		}
	}
	return ta.Widget[0]
}

func (ta *App) NextWidgets(current tview.Primitive) tview.Primitive {
	for i, w := range ta.Widget {
		if w == current {
			if i+1 < len(ta.Widget) {
				return ta.Widget[i+1]
			}
			return ta.Widget[0]
		}
	}
	return ta.Widget[0]
}
