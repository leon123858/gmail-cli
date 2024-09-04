package dashboard

import (
	"github.com/leon123858/gmail-cli/tui"
	"github.com/rivo/tview"
)

var app *tui.App

type PageName string

const (
	MailsBoard PageName = "MailsBoard"
)

func init() {
	app = tui.NewApp()
}

func GetRootPages() *tview.Pages {
	return app.GetPage()
}

func ShowBoard(name PageName) {
	switch name {
	case MailsBoard:
		app.AddPage(string(MailsBoard), NewMailsBoard())
	}
	app.ShowPage(string(name))
}
