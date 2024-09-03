package tui

import "github.com/rivo/tview"

type List struct {
	*tview.List
}

func NewList() *List {
	return &List{
		List: tview.NewList(),
	}
}

func (l *List) SetSelectedFunc(f func(int, string, string, rune)) {
	l.List.SetSelectedFunc(f)
}

func (l *List) SetChangedFunc(f func(int, string, string, rune)) {
	l.List.SetChangedFunc(f)
}

func (l *List) AddItem(text, secondaryText string, shortcut rune, selected func()) {
	l.List.AddItem(text, secondaryText, shortcut, selected)
}
