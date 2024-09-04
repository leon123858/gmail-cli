package dashboard

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/leon123858/gmail-cli/gmail"
	"github.com/leon123858/gmail-cli/tui"
	"github.com/spf13/viper"
)

func NewMailsBoard() *tui.Flex {
	leftList := tui.NewList()
	rightText := tui.NewText()

	leftList.SetBorder(true).SetTitle("Emails")
	rightText.SetBorder(true).SetTitle("Text")

	// init mail list
	accounts := viper.GetStringSlice("accounts")
	if len(accounts) == 0 {
		fmt.Println("No accounts configured. Use 'gmail-cli config add <email>' to add an account.")
	}

	flex := tui.NewFlex(false)

	flex.AddItem(leftList, 0, 1, false)
	flex.AddItem(rightText, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			app.SetFocus(rightText)
			return nil
		case tcell.KeyLeft:
			app.SetFocus(leftList)
			return nil
		default:
			return event
		}
	})

	{
		var mailResChans []chan gmail.MailResChanel
		pageSize := 30
		for i, account := range accounts {
			mailResChans = append(mailResChans, make(chan gmail.MailResChanel))
			go gmail.ReadEmails(account, pageSize, mailResChans[i])
		}

		for _, mailResChan := range mailResChans {
			go func(mailResChan chan gmail.MailResChanel) {
				for {
					res := <-mailResChan
					if res.Err != nil {
						if res.Err.Error() == "EOF" {
							break
						}
						continue
					}
					leftList.AddItem(res.Res.Subject, res.Account, 0, func() {
						// set right text, as mail body
						rightText.SetText(fmt.Sprintf("From: %s\nDate: %s\n\n%s", res.Res.From, res.Res.Date, res.Res.Body), false)
					})
				}
			}(mailResChan)
		}
	}

	app.SetFocus(leftList)

	return flex
}
