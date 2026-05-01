# Mastodon Usage

Ensure you have a Mastodon account and have created an application
to obtain your client ID, client secret, and access token.

## Sample Code

```go
package main

import (
    "context"
    "log"

    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/mastodon"
)

func main() {
    mastodonSvc := mastodon.New(
        "https://mastodon.social",  // Your instance URL
        "your-client-id",
        "your-client-secret",
        "your-access-token",
    )

    mastodonSvc.AddReceivers("@user@mastodon.social")

    notifier := notify.New()
    notifier.UseServices(mastodonSvc)

    err := notifier.Send(context.Background(), "Subject", "Message body")
    if err != nil {
        log.Fatal(err)
    }
}
```
