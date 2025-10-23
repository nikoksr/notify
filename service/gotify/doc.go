/*
Package gotify provides message notification integration for Gotify (https://gotify.net).

Gotify is a simple server for sending and receiving messages in real-time per WebSocket.
It's perfect for sending notifications from your applications.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/gotify"
	)

	func main() {
		gotifyService := gotify.New("https://gotify.example.com", "your-app-token")

		notify.UseServices(gotifyService)

		err := notify.Send(context.Background(), "Subject", "Message")
		if err != nil {
			log.Fatal(err)
		}
	}
*/
package gotify
