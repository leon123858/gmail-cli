package cli

import (
	"context"
	"fmt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"strings"
	"sync"
	"time"
)

func readEmails(account string, numEmails int, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := getClient(account)
	if err != nil {
		log.Printf("Failed to get client for %s: %v", account, err)
		return
	}

	gmailService, err := gmail.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Failed to create Gmail service for %s: %v", account, err)
		return
	}

	// 獲取今天的日期，格式為 YYYY/MM/DD
	today := time.Now().Add(-24 * time.Hour).Format("2006/01/02")
	query := fmt.Sprintf("after:%s", today)

	msgs, err := gmailService.Users.Messages.List("me").Q(query).MaxResults(int64(numEmails)).Do()
	if err != nil {
		log.Printf("Failed to retrieve messages for %s: %v", account, err)
		return
	}

	for _, msg := range msgs.Messages {
		m, err := gmailService.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			log.Printf("Failed to retrieve message details for %s: %v", account, err)
			continue
		}

		var subject, _, _ string
		for _, header := range m.Payload.Headers {
			switch header.Name {
			case "Subject":
				subject = header.Value
				//case "From":
				//	from = header.Value
				//case "Date":
				//	date = header.Value
			}
		}

		fmt.Printf("%s: %s\n", account, subject)
		fmt.Println(strings.Repeat("-", 40))
	}
}
