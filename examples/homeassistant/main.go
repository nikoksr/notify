package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/homeassistant"
)

func main() {
	// Create Home Assistant service with your instance URL and access token
	haService := homeassistant.New("http://homeassistant.local:8123", "your-long-lived-access-token-here")

	// Use the service with notify
	notify.UseServices(haService)

	// Send a notification
	err := notify.Send(
		context.Background(),
		"Hello from notify!",
		"This is a test message sent via Home Assistant service. Check your Home Assistant notifications!",
	)
	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully!")
	log.Println("Check your Home Assistant web interface or mobile app to see the notification.")
}