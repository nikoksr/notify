/*
Package mastodon provides message notification integration for Mastodon.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/mastodon"
	)

	func main() {
		mastodonSvc := mastodon.New(
			"https://mastodon.social",
			"your-client-id",
			"your-client-secret",
			"your-access-token",
		)

		mastodonSvc.AddReceivers("@user@mastodon.social")

		notifier := notify.New()
		notifier.UseServices(mastodonSvc)

		if err := notifier.Send(context.Background(), "Subject", "Message"); err != nil {
			log.Fatal(err)
		}
	}
*/
package mastodon
