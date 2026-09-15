# Zulip

## Steps for creating a Zulip Bot

Follow the below instructions to create a bot email and api key required for the service:

1. Create a Zulip Organization
2. Go to settings and create a new bot. Copy the bot email and api key.
3. Copy your Oranization URL. Copy the entire url `https://your-domain.zulipchat.com`
4. Copy the stream name of the stream and its topic if you want to post a message to stream or just copy an email address of the receiver.

## Sample Code

```go
	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/zulip"
	)

	func main() {
	    zulipSvc, err := zulip.New("server-base-url", "bot-email", "api-key")
		if err != nil {
			log.Fatalf("zulip.New() failed: %s", err.Error())
		}

	    zulipSvc.AddReceivers(Direct("test2@gmail.com"), Stream("stream", "topic"))

		notifier := notify.New()
		notifier.UseServices(zulipSvc)

		err = notifier.Send(context.Background(), "subject", "message")
		if err != nil {
			log.Fatalf("notifier.Send() failed: %s", err.Error())
		}

		log.Println("notification sent")
	}
```
