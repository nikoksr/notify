# Zulip Usage

Ensure that you have already navigated to your GOPATH and installed the following packages:

* `go get -u github.com/nikoksr/notify`

## Steps for creating Zulip Bot

These are general and very high level instructions

1. Create a new Zulip bot (https://zulip.com/help/add-a-bot-or-integration)
2. Copy your *Organization URL* from the browser address bar. You need to copy only subdomain `your-org` from the full url `your-org.zulipchat.com` without the hostname `.zulipchat.com`.
3. Copy your *Bot Email* and *API Key* for usage below
4. Copy the *Stream name* of the stream if you want to post a message to stream or just copy an email address of the receiver.
5. Now you should be good to use the code below

## Sample Code

```go
package main

import (
    "context"
    "fmt"
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/zulip"
)

func main() {

    notifier := notify.New()

    // Provide your Zulip Bot credentials
    zulipService := zulip.New(
        "your-org",
        "ZULIP_API_KEY",
        "email-bot@your-org.zulipchat.com",
    )

    // Passing a Zulip receivers as a receiver for our messages.
    // Where to send our messages.
    // It can be direct or stream message
    zulipService.AddReceivers(zulip.Direct("some-user@email.com"))
    zulipService.AddReceivers(zulip.Stream("alerts", "critical"))

    // Tell our notifier to use the Zulip service. You can repeat the above process
    // for as many services as you like and just tell the notifier to use them.
    notifier.UseServices(zulipService)

    // Send a message
    err := notifier.Send(
        context.Background(),
        "Hello from notify :wave:\n",
        "Message written in Go!",
    )

    if err != nil {
        fmt.Println(err)
    }

}
```
