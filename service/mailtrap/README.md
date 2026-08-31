# Mailtrap

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/mailtrap)

[Mailtrap](https://mailtrap.io) is an email delivery platform for sending, testing, and controlling emails. This service sends transactional emails through the Mailtrap [Email Sending API](https://api-docs.mailtrap.io).

## Prerequisites

To use the Mailtrap notification service, you will need:

- A Mailtrap API token — create one under [API Tokens](https://mailtrap.io/api-tokens).
- A verified sending domain and a sender email address that belongs to it.

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mailtrap"
)

func main() {
	mailtrapSvc := mailtrap.New("your-api-token", "sender@your-domain.com")

	mailtrapSvc.AddReceivers("recipient@example.com")

	notifier := notify.New()
	notifier.UseServices(mailtrapSvc)

	if err := notifier.Send(context.Background(), "subject", "message"); err != nil {
		log.Fatalf("notifier.Send() failed: %s", err.Error())
	}

	log.Println("notification sent")
}
```

## Options

### Sender name

By default only the sender email address is sent. Use `WithSenderName` to also set a display name:

```go
mailtrapSvc := mailtrap.New(
	"your-api-token",
	"sender@your-domain.com",
	mailtrap.WithSenderName("Acme Notifications"),
)
```

### Plain-text messages

The message body is sent as HTML by default. Switch to plain text with `BodyFormat`:

```go
mailtrapSvc := mailtrap.New("your-api-token", "sender@your-domain.com")
mailtrapSvc.BodyFormat(mailtrap.PlainText)
```

### Bulk sending

To route messages through the Mailtrap bulk sending host (intended for high-volume sending), use `WithBulkHost`:

```go
mailtrapSvc := mailtrap.New(
	"your-api-token",
	"sender@your-domain.com",
	mailtrap.WithBulkHost(),
)
```

### Sandbox (testing)

To capture emails in a Mailtrap test inbox instead of delivering them, use `WithSandbox` with your inbox ID. The sandbox API requires the inbox ID as part of the request path:

```go
mailtrapSvc := mailtrap.New(
	"your-api-token",
	"sender@your-domain.com",
	mailtrap.WithSandbox("your-inbox-id"),
)
```

### Custom HTTP client

To configure timeouts, proxies, or a custom transport, provide your own HTTP client with `WithHTTPClient`:

```go
import (
	"net/http"
	"time"

	"github.com/nikoksr/notify/service/mailtrap"
)

httpClient := &http.Client{Timeout: 10 * time.Second}

mailtrapSvc := mailtrap.New(
	"your-api-token",
	"sender@your-domain.com",
	mailtrap.WithHTTPClient(httpClient),
)
```
