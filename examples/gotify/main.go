package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/gotify"
)

func main() {
	// Create gotify service with your server URL and app token
	gotifyService := gotify.New("https://gotify.example.com", "your-app-token-here")

	// Use the service with notify
	notify.UseServices(gotifyService)

	// Send a notification
	err := notify.Send(
		context.Background(),
		"Hello from notify!",
		"This is a test message sent via Gotify service. Check your Gotify app to see this notification!",
	)
	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully!")
	log.Println("Check your Gotify app or web interface to see the message.")
}