/*
Package ntfy provides message notification integration for ntfy (https://ntfy.sh).

Ntfy is a simple HTTP-based pub-sub notification service. It allows you to send
notifications via scripts from any computer, entirely without signup, cost or setup.
You can also self-host your own ntfy server.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/ntfy"
	)

	func main() {
		ntfyService := ntfy.New("ntfy.sh")
		ntfyService.AddReceivers("my-topic")

		notify.UseServices(ntfyService)

		err := notify.Send(context.Background(), "Subject", "Message")
		if err != nil {
			log.Fatal(err)
		}
	}
*/
package ntfy