# Nylas

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/nylas)

## Prerequisites

To use the Nylas notification service, you will need:

- **Nylas API Key** - Your application's API key
- **Grant ID** - A grant representing an authenticated email account
- **Email Address** - The sender email address associated with the grant

## Getting Started with Nylas

### 1. Sign Up for Nylas

Visit the [Nylas v3 Dashboard](https://dashboard-v3.nylas.com/) and create a free account.

### 2. Create an Application

1. After signing in, create a new application in the dashboard
2. Note your **API Key** from the application settings

### 3. Connect an Email Account (Create a Grant)

A **Grant** in Nylas v3 represents an authenticated connection between your application and a user's email provider (Gmail, Outlook, Exchange, etc.).

**To create a grant:**

1. In your Nylas dashboard, navigate to **Grants** or **Connected Accounts**
2. Click **"Add Grant"** or **"Connect Account"**
3. Choose your email provider (Google, Microsoft, etc.)
4. Follow the OAuth flow to authenticate your email account
5. Once connected, copy the **Grant ID** - it looks like: `12d6d3d7-2441-4083-ab45-cc1525edd1f7`

**Important:** Each grant is tied to a specific email account. To send emails, you need a grant for the email address you want to send from.

### 4. Find Your Credentials

From the Nylas dashboard, you should now have:
- ✅ **API Key** (from Application settings)
- ✅ **Grant ID** (from the connected account/grant)
- ✅ **Email Address** (the email you connected)

## Understanding Nylas v3 Authentication

Nylas v3 uses a **Grant-based authentication model**:

- **Grant ID**: A UUID representing an authenticated email account connection
- **API Key**: Your application-level credential (kept secret)
- All API requests require both the API Key (in headers) and Grant ID (in the URL path)

This is different from v2, which used account IDs and access tokens.

## Usage

### Basic Example

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/nylas"
)

func main() {
	// Load credentials from environment variables (recommended)
	apiKey := os.Getenv("NYLAS_API_KEY")
	grantID := os.Getenv("NYLAS_GRANT_ID")
	senderEmail := os.Getenv("NYLAS_SENDER_EMAIL")

	// Create Nylas service
	nylasService := nylas.New(
		apiKey,
		grantID,
		senderEmail,
		"Your Name", // Sender display name
	)

	// Add recipients
	nylasService.AddReceivers("[email protected]", "[email protected]")

	// Create notifier and add the service
	notifier := notify.New()
	notifier.UseServices(nylasService)

	// Send notification
	err := notifier.Send(
		context.Background(),
		"Welcome Email",
		"<h1>Hello!</h1><p>This is a test email via Nylas.</p>",
	)
	if err != nil {
		log.Fatalf("Failed to send: %v", err)
	}

	log.Println("Email sent successfully!")
}
```

### HTML vs Plain Text

By default, emails are sent as HTML. To send plain text:

```go
nylasService := nylas.New(apiKey, grantID, senderEmail, "Sender Name")
nylasService.BodyFormat(nylas.PlainText)
```

### Regional Configuration

Nylas has different API endpoints for different regions. By default, the US region is used.

**For EU region (recommended approach):**

```go
// Use NewWithRegion for cleaner, type-safe region selection
nylasService := nylas.NewWithRegion(
	apiKey,
	grantID,
	senderEmail,
	"Sender Name",
	nylas.RegionEU,
)
```

**Available regions:**
- `nylas.RegionUS` - United States (default)
- `nylas.RegionEU` - European Union

**Alternative approach (manual URL):**

```go
// You can also set a custom base URL manually
nylasService := nylas.New(apiKey, grantID, senderEmail, "Sender Name")
nylasService.WithBaseURL("https://api.eu.nylas.com")
```

### Custom HTTP Client

For advanced use cases (custom timeouts, proxies, etc.):

```go
import (
	"net/http"
	"time"
)

customClient := &http.Client{
	Timeout: 60 * time.Second,
}

nylasService := nylas.New(apiKey, grantID, senderEmail, "Sender Name")
nylasService.WithHTTPClient(customClient)
```

### Complete Example

```go
package main

import (
	"context"
	"log"
	"os"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/nylas"
)

func main() {
	// Determine region
	region := nylas.RegionUS
	if os.Getenv("NYLAS_REGION") == "EU" {
		region = nylas.RegionEU
	}

	// Create Nylas service with credentials and region
	nylasService := nylas.NewWithRegion(
		os.Getenv("NYLAS_API_KEY"),
		os.Getenv("NYLAS_GRANT_ID"),
		os.Getenv("NYLAS_SENDER_EMAIL"),
		"Notification Service",
		region,
	)

	// Set email format
	nylasService.BodyFormat(nylas.HTML)

	// Add multiple recipients
	nylasService.AddReceivers(
		"[email protected]",
		"[email protected]",
		"[email protected]",
	)

	// Use with notify
	notifier := notify.New()
	notifier.UseServices(nylasService)

	// Send rich HTML email
	ctx := context.Background()
	err := notifier.Send(ctx, "System Alert", `
		<html>
		<body>
			<h2>⚠️ Important Notification</h2>
			<p>Your system requires attention.</p>
			<ul>
				<li>Status: <strong>Warning</strong></li>
				<li>Time: <strong>2026-01-17 10:00:00</strong></li>
			</ul>
		</body>
		</html>
	`)

	if err != nil {
		log.Fatalf("Failed to send notification: %v", err)
	}

	log.Println("Notification sent successfully!")
}
```

## Environment Variables

For security best practices, store your credentials in environment variables:

```bash
export NYLAS_API_KEY="nyk_v0_..."
export NYLAS_GRANT_ID="12d6d3d7-2441-4083-ab45-cc1525edd1f7"
export NYLAS_SENDER_EMAIL="[email protected]"
export NYLAS_REGION="US"  # or "EU"
```

## Troubleshooting

### "No receivers configured" error
Make sure you call `AddReceivers()` before sending:
```go
nylasService.AddReceivers("[email protected]")
```

### "Nylas API error: Invalid grant ID"
- Verify your Grant ID is correct (check the Nylas dashboard)
- Ensure the grant is still active (not revoked)
- Confirm you're using the correct API key for your application

### Timeout errors
For self-hosted Exchange servers, emails can take up to 150 seconds to send. The default timeout is configured for this, but you can adjust it:
```go
customClient := &http.Client{Timeout: 180 * time.Second}
nylasService.WithHTTPClient(customClient)
```

## Resources

- [Nylas v3 Dashboard](https://dashboard-v3.nylas.com/)
- [Nylas API Documentation](https://developer.nylas.com/docs/v3/)
- [Sending Email Guide](https://developer.nylas.com/docs/v3/email/send-email/)
- [Authentication & Grants](https://developer.nylas.com/docs/v3/auth/)
- [API Reference](https://developer.nylas.com/docs/v3/api-references/)

## Implementation Notes

This service uses **direct REST API calls** to Nylas v3 endpoints, as there is no official Go SDK for Nylas v3. The implementation:

- ✓ Has zero external dependencies (beyond Go standard library)
- ✓ Follows Nylas v3 API specifications exactly
- ✓ Includes comprehensive error handling with detailed error messages
- ✓ Supports context-based cancellation
- ✓ Production-ready with appropriate timeouts
