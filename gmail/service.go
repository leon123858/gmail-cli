package gmail

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"log"
	"strings"
	"time"
)

func removeHTMLTagsWithGoquery(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	// 獲取所有文本節點
	doc.Find("style").Remove() // 移除所有 style 標籤
	text := doc.Find("*").Text()
	// 去除空白字符
	text = strings.TrimSpace(text)
	return text
}

func base64Decode(s string) string {
	data, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return "error decoding base64"
	}
	return string(data)
}

type Mail struct {
	Subject string
	From    string
	Date    string
	Body    string
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
		var body string
		switch m.Payload.MimeType {
		case "text/plain":
			body = base64Decode(m.Payload.Body.Data)
		case "text/html":
			body = removeHTMLTagsWithGoquery(base64Decode(m.Payload.Body.Data))
		case "multipart/alternative":
			for _, part := range m.Payload.Parts {
				switch part.MimeType {
				case "text/plain":
					body += base64Decode(part.Body.Data)
				case "text/html":
					body += removeHTMLTagsWithGoquery(base64Decode(part.Body.Data))
				default:
					body += "Unknown message type: " + part.MimeType
				}
			}
		case "multipart/mixed":
			for _, part := range m.Payload.Parts {
				switch part.MimeType {
				case "text/plain":
					body += base64Decode(part.Body.Data)
				case "text/html":
					body += removeHTMLTagsWithGoquery(base64Decode(part.Body.Data))
				default:
					body += "Unknown message type: " + part.MimeType
				}
			}
		default:
			body = "Unknown message type: " + m.Payload.MimeType
		}

		ch <- MailResChanel{Res: Mail{Subject: subject, From: from, Date: date, Body: body}, Account: account}

		//fmt.Printf("%s: %s\n", account, subject)
		//fmt.Println(strings.Repeat("-", 40))
	}

	ch <- MailResChanel{Err: errors.New("EOF")}
}
