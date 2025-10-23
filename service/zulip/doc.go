/*
Package zulip provides message notification integration for Zulip (https://zulip.com).

Zulip is an open-source team chat application with powerful search and integrations.
This service supports both stream messages and direct messages.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/zulip"
	)

	func main() {
		zulipService := zulip.New("https://company.zulipchat.com", "bot@company.com", "api-key")

		notify.UseServices(zulipService)

		err := notify.Send(context.Background(), "Subject", "Message")
		if err != nil {
			log.Fatal(err)
		}
	}
*/
package zulip
