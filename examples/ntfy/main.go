package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/ntfy"
)

func main() {
	// Create ntfy service using public ntfy.sh server
	ntfyService := ntfy.New("ntfy.sh")

	// Add topics - these are like channels you subscribe to
	ntfyService.AddReceivers("notify-example-topic")

	// Use the service with notify
	notify.UseServices(ntfyService)

	// Send a notification
	err := notify.Send(
		context.Background(),
		"Hello from notify!",
		"This is a test message sent via ntfy service. Subscribe to 'notify-example-topic' on ntfy.sh to see this message!",
	)
	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully!")
	log.Println("To receive this message, visit: https://ntfy.sh/notify-example-topic")
}