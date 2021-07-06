# Twilio

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/twilio)

## Prerequisites

You will need a [Twilio](https://www.twilio.com/) account for SID and AuthToken.

```go
package main

import (
        "context"
        "log"

        "github.com/nikoksr/notify"
        "github.com/nikoksr/notify/service/twilio"
)

func main() {

       twilioSvc, err := twilio.New("YOUR_PHONE_NUMBER", "YOUR_TWILIO_SID", "YOUR_TWILIO_TOKEN")
        if err != nil {
                log.Fatalf("twilio.New() failed: %s", err.Error())
        }

        twilioSvc.AddReceivers("Contact1", "Contact2", "Contact3,.......")

        notifier := notify.New()
        notifier.UseServices(twilioSvc)

        err = notifier.Send(context.Background(), "subject", "message")
        if err != nil {
                log.Fatalf("notifier.Send() failed: %s", err.Error())
        }

        log.Println("notification sent")
}
```
