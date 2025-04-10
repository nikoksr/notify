/*
Package mastodon provides message notification integration for the Mastodon social network.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/mastodon"
	)

	func main() {
		cfg := &mastodon.Config{
			ServerURL:   "https://mastodon.social",
			AccessToken: "your-access-token",
			Visibility:  mastodon.VisibilityPublic,
		}

		service, err := mastodon.New(cfg)
		if err != nil {
			log.Fatalf("mastodon.New() failed: %s", err.Error())
		}

		notifier := notify.New()
		notifier.UseServices(service)

		err = notifier.Send(context.Background(), "", "Hello from Mastodon!")
		if err != nil {
			log.Fatalf("notifier.Send() failed: %s", err.Error())
		}
	}

Visibility values:

  - VisibilityPublic   — visible to everyone
  - VisibilityUnlisted — visible to everyone but not listed
  - VisibilityPrivate  — visible to followers only
  - VisibilityDirect   — visible only to mentioned users (direct message)
*/
package mastodon
