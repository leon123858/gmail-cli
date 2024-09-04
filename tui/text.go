package tui

import "github.com/rivo/tview"

type Text struct {
	*tview.TextArea
}

func NewText() *Text {
	return &Text{
		TextArea: tview.NewTextArea(),
	}
}

func (t *Text) SetBorder(show bool) *Text {
	t.TextArea.SetBorder(show)
	return t
}

func (t *Text) SetTitle(title string) *Text {
	t.TextArea.SetTitle(title)
	return t
}

func (t *Text) SetText(text string, cursorAtEnd bool) *Text {
	t.TextArea.SetText(text, cursorAtEnd)
	return t
}

func (t *Text) SetWordWrap(wrap bool) *Text {
	t.TextArea.SetWordWrap(wrap)
	return t
}
