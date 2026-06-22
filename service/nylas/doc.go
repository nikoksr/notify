/*
Package nylas provides a service for sending email notifications via Nylas API v3.

Nylas is a communications platform that provides APIs for email, calendar, and contacts.
This service implements support for sending emails through Nylas API v3 using direct
REST API calls.

Usage:

	package main

	import (
	    "context"
	    "log"

	    "github.com/nikoksr/notify"
	    "github.com/nikoksr/notify/service/nylas"
	)

	func main() {
	    // Create a Nylas service with your API credentials.
	    // You'll need:
	    // - API Key: Your Nylas application API key
	    // - Grant ID: The grant ID for the email account to send from
	    // - Sender Address: The email address to send from
	    // - Sender Name: The display name for the sender (optional)
	    nylasService := nylas.New(
	        "your_api_key",
	        "your_grant_id",
	        "[email protected]",
	        "Your Name",
	    )

	    // Add email addresses to send to.
	    nylasService.AddReceivers("[email protected]", "[email protected]")

	    // Optional: Set body format (default is HTML).
	    nylasService.BodyFormat(nylas.HTML)

	    // Tell our notifier to use the Nylas service.
	    notify.UseServices(nylasService)

	    // Send a test message.
	    err := notify.Send(
	        context.Background(),
	        "Test Subject",
	        "<h1>Hello!</h1><p>This is a test message from Nylas.</p>",
	    )
	    if err != nil {
	        log.Fatalf("Failed to send notification: %v", err)
	    }
	}

Regional Configuration:

For EU region or other Nylas API regions, use NewWithRegion():

	// Create a Nylas service for the EU region
	nylasService := nylas.NewWithRegion(
	    "your_api_key",
	    "your_grant_id",
	    "[email protected]",
	    "Your Name",
	    nylas.RegionEU,
	)

Supported regions: nylas.RegionUS (default), nylas.RegionEU

For more information about Nylas API v3, see:
  - Getting Started: https://developer.nylas.com/docs/v3/getting-started/
  - Sending Email: https://developer.nylas.com/docs/v3/email/send-email/
  - API Reference: https://developer.nylas.com/docs/v3/api-references/
*/
package nylas
