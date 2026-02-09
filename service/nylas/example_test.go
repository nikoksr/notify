package nylas_test

import (
	"context"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/nylas"
)

func Example() {
	// Create a Nylas service instance with your credentials.
	// Note: In production, use environment variables or a secure config for credentials.
	nylasService := nylas.New(
		"your-api-key",      // Nylas API key
		"your-grant-id",     // Grant ID for the email account
		"[email protected]", // Sender email address
		"Your Name",         // Sender display name
	)

	// Add one or more recipient email addresses.
	nylasService.AddReceivers("[email protected]", "[email protected]")

	// Optional: Set the body format (default is HTML).
	nylasService.BodyFormat(nylas.HTML)

	// Optional: Use a different region (e.g., EU region).
	// nylasService.WithBaseURL("https://api.eu.nylas.com")

	// Create a notifier and add the Nylas service.
	notifier := notify.New()
	notifier.UseServices(nylasService)

	// Send a notification.
	_ = notifier.Send(
		context.Background(),
		"Welcome to Notify with Nylas!",
		"<h1>Hello!</h1><p>This is an email notification sent via Nylas API v3.</p>",
	)

	// Output:
}
