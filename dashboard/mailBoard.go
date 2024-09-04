package dashboard

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/leon123858/gmail-cli/gmail"
	"github.com/leon123858/gmail-cli/tui"
	"github.com/spf13/viper"
	"sync"
)

func NewMailsBoard() *tui.Flex {
	leftList := tui.NewList()
	rightList := tui.NewList()

	leftList.SetBorder(true).SetTitle("Emails")
	rightList.SetBorder(true).SetTitle("Text")

	// init mail list
	accounts := viper.GetStringSlice("accounts")
	if len(accounts) == 0 {
		fmt.Println("No accounts configured. Use 'gmail-cli config add <email>' to add an account.")
	}

	flex := tui.NewFlex(false)

	flex.AddItem(leftList, 0, 1, false)
	flex.AddItem(rightList, 0, 1, false)

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRight:
			app.SetFocus(rightList)
			return nil
		case tcell.KeyLeft:
			app.SetFocus(leftList)
			return nil
		default:
			return event
		}
	})

	var mailResChans []chan gmail.MailResChanel
	pageSize := 10
	for i, account := range accounts {
		mailResChans = append(mailResChans, make(chan gmail.MailResChanel))
		go gmail.ReadEmails(account, pageSize, mailResChans[i])
	}

	wg := sync.WaitGroup{}
	wg.Add(len(mailResChans))

	for _, mailResChan := range mailResChans {
		go func(mailResChan chan gmail.MailResChanel) {
			for {
				res := <-mailResChan
				if res.Err != nil {
					if res.Err.Error() == "EOF" {
						wg.Done()
						break
					}
					continue
				}
				leftList.AddItem(res.Res.Subject, res.Account, 0, nil)
			}
		}(mailResChan)
	}
	wg.Wait()

	app.SetFocus(leftList)

	return flex
}
