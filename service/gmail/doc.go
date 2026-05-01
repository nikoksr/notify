/*
Package gmail implements a Gmail notifier, allowing plain text emails to be sent
to multiple recipients through the Gmail API using OAuth2 credentials.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/gmail"
		"golang.org/x/oauth2"
	)

	func main() {
		notifier := notify.New()

		token := &oauth2.Token{AccessToken: "ACCESS_TOKEN"}
		tokenSource := oauth2.StaticTokenSource(token)

		gmailService, err := gmail.New(tokenSource, "sender@example.com")
		if err != nil {
			log.Fatal(err)
		}

		gmailService.AddReceivers("alice@example.com", "bob@example.com")

		notifier.UseServices(gmailService)

		_ = notifier.Send(
			context.Background(),
			"Hello!",
			"I am a bot written in Go!",
		)
	}
*/
package gmail
