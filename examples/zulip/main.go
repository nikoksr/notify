package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/zulip"
)

func main() {
	// Create Zulip service with your server URL, bot email, and API key
	zulipService := zulip.New("https://company.zulipchat.com", "bot@company.com", "your-api-key-here")

	// Use the service with notify
	notify.UseServices(zulipService)

	// Send a notification to default stream
	err := notify.Send(
		context.Background(),
		"Hello from notify!",
		"This is a test message sent via Zulip service. Check your Zulip streams!",
	)
	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	// Send to specific stream with topic
	err = zulipService.SendToStream(
		context.Background(),
		"development",
		"Notifications",
		"Testing the new notify integration!",
	)
	if err != nil {
		log.Fatalf("Failed to send stream message: %v", err)
	}

	// Send direct message
	err = zulipService.SendDirectMessage(
		context.Background(),
		"user@company.com",
		"Private notification test from notify library!",
	)
	if err != nil {
		log.Fatalf("Failed to send direct message: %v", err)
	}

	log.Println("All notifications sent successfully!")
	log.Println("Check your Zulip streams and direct messages.")
}
