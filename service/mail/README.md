# Mail

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/twilio)

## Prerequisites

1. Enable the SMTP service for the sender's mail
2. You will need to following information to be able to send mail.

- SMTP server host: e.g. smtp.163.com
- SMTP server port: e.g. 25, 465, 587
- username: it's sender address if no special instructions from the mail provider
- password: authorization password
- sender address: e.g. xxx@163.com
- receiver addresses: e.g. xxx@163.com

## Usage

```go
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mail"
)

func main() {
	smtpHost := "smtp.163.com"
	smtpPort := 465
	username := "sender@163.com"
	password := "xxx"
	senderAddress := "sender@163.com"
	receiverAddresses := []string{"receiver1@163.com", "receiver2@gmail.com"}

	notifier := notify.New()

	smtpHostAddress := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	mailService := mail.New(senderAddress, smtpHostAddress)
	mailService.AuthenticateSMTP("", username, password, smtpHost)
	mailService.AddReceivers(receiverAddresses...)

	// NOTICE: If your smtpPort is 25, just comment out this line of code,
	// Otherwise notifier.Send() will report an error:
	// "failed to send mail: tls: first record does not look like a TLS handshake: send notification"
	mailService.SetTLS(&tls.Config{ServerName: smtpHost})

	notifier.UseServices(mailService)

	err := notifier.Send(context.Background(), "subject", "message")
	if err != nil {
		log.Fatalf("notifier.Send() failed: %s", err.Error())
	}

	log.Println("notification sent")
}
```
