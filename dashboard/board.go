package dashboard

import (
	"github.com/leon123858/gmail-cli/tui"
	"github.com/rivo/tview"
)

var app *tui.App

type PageName string

const (
	Board PageName = "board"
)

func init() {
	app = tui.NewApp()
	leftList := tui.NewList()
	rightList := tui.NewList()

	leftList.AddItem("Inbox", "", 'i', nil)
	leftList.AddItem("Sent", "", 's', nil)
	leftList.AddItem("Drafts", "", 'd', nil)
	leftList.AddItem("Trash", "", 't', nil)

	rightList.AddItem("Email 1", "", '1', nil)
	rightList.AddItem("Email 2", "", '2', nil)
	rightList.AddItem("Email 3", "", '3', nil)

	flex := tui.NewFlex(false)

	flex.AddItem(leftList, 10, 1, true)
	flex.AddItem(rightList, 10, 1, false)

	app.AddPage(string(Board), flex)
}

func GetRootPages() *tview.Pages {
	return app.GetPage()
}

func ShowBoard(name PageName) {
	app.ShowPage(string(name))
}
