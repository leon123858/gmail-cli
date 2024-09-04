package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Flex struct {
	*tview.Flex
}

func NewFlex(isRow bool) *Flex {
	if isRow {
		return &Flex{
			Flex: tview.NewFlex().SetDirection(tview.FlexRow),
		}
	}
	return &Flex{
		Flex: tview.NewFlex().SetDirection(tview.FlexColumn),
	}
}

func (f *Flex) AddItem(item tview.Primitive, size int, proportion int, focus bool) *Flex {
	f.Flex.AddItem(item, size, proportion, focus)
	return f
}

func (f *Flex) SetInputCapture(capture func(event *tcell.EventKey) *tcell.EventKey) *Flex {
	f.Flex.SetInputCapture(capture)
	return f
}

func (f *Flex) SetBorder(show bool) *Flex {
	f.Flex.SetBorder(show)
	return f
}
