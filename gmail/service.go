package gmail

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"time"
)

type Mail struct {
	Subject string
	From    string
	Date    string
}

type MailResChanel struct {
	Account string
	Res     Mail
	Err     error
}

func ReadEmails(account string, numEmails int, ch chan MailResChanel) {
	client, err := getClient(account)
	if err != nil {
		log.Printf("Failed to get client for %s: %v", account, err)
		ch <- MailResChanel{Err: errors.New("EOF")}
		return
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Failed to create Gmail service for %s: %v", account, err)
		ch <- MailResChanel{Err: errors.New("EOF")}
		return
	}

	// 獲取今天的日期，格式為 YYYY/MM/DD
	today := time.Now().Add(-24 * time.Hour).Format("2006/01/02")
	query := fmt.Sprintf("after:%s", today)

	msgs, err := gmailService.Users.Messages.List("me").Q(query).MaxResults(int64(numEmails)).Do()
	if err != nil {
		log.Printf("Failed to retrieve messages for %s: %v", account, err)
		ch <- MailResChanel{Err: errors.New("EOF")}
		return
	}

	for _, msg := range msgs.Messages {
		m, err := gmailService.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			log.Printf("Failed to retrieve message details for %s: %v", account, err)
			ch <- MailResChanel{Err: errors.New("failed to retrieve message details")}
			continue
		}

		var subject, from, date string
		for _, header := range m.Payload.Headers {
			switch header.Name {
			case "Subject":
				subject = header.Value
			case "From":
				from = header.Value
			case "Date":
				date = header.Value
			}
		}

		ch <- MailResChanel{Res: Mail{Subject: subject, From: from, Date: date}, Account: account}

		//fmt.Printf("%s: %s\n", account, subject)
		//fmt.Println(strings.Repeat("-", 40))
	}

	ch <- MailResChanel{Err: errors.New("EOF")}
}
