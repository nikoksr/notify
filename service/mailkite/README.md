# MailKite

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/mailkite)

[MailKite](https://mailkite.dev) is an email API for sending and receiving transactional mail. This service sends notifications through the MailKite [Send API](https://mailkite.dev/docs/sending) using only the Go standard library.

## Prerequisites

To use the MailKite notification service, you will need:

- A MailKite API key (`mk_live_…`) — create one in the [dashboard](https://app.mailkite.dev).
- A sender email address on a domain verified in that account.

## Usage

```go
package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/mailkite"
)

func main() {
	mailkiteSvc := mailkite.New("mk_live_your-api-key", "sender@your-domain.com")

	mailkiteSvc.AddReceivers("recipient@example.com")

	notifier := notify.New()
	notifier.UseServices(mailkiteSvc)

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
mailkiteSvc := mailkite.New(
	"mk_live_your-api-key",
	"sender@your-domain.com",
	mailkite.WithSenderName("Acme Notifications"),
)
```

### Reply-To

To have replies land somewhere other than the sender address, use `WithReplyTo`:

```go
mailkiteSvc := mailkite.New(
	"mk_live_your-api-key",
	"sender@your-domain.com",
	mailkite.WithReplyTo("support@your-domain.com"),
)
```

### Plain-text messages

The message body is sent as HTML by default. Switch to plain text with `BodyFormat`:

```go
mailkiteSvc := mailkite.New("mk_live_your-api-key", "sender@your-domain.com")
mailkiteSvc.BodyFormat(mailkite.PlainText)
```

### Custom base URL

To point the client at a different API endpoint (for example a local mock in tests), use `WithBaseURL`:

```go
mailkiteSvc := mailkite.New(
	"mk_live_your-api-key",
	"sender@your-domain.com",
	mailkite.WithBaseURL("http://localhost:8787"),
)
```

### Custom HTTP client

To configure timeouts, proxies, or a custom transport, provide your own HTTP client with `WithHTTPClient`:

```go
import (
	"net/http"
	"time"

	"github.com/nikoksr/notify/service/mailkite"
)

httpClient := &http.Client{Timeout: 10 * time.Second}

mailkiteSvc := mailkite.New(
	"mk_live_your-api-key",
	"sender@your-domain.com",
	mailkite.WithHTTPClient(httpClient),
)
```
