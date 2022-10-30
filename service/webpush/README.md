# Webpush Notifications

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/webpush)


## Prerequisites

Generate VAPID Public and Private Keys for the notification service. This can be done using many tools, one of which is [`GenerateVAPIDKeys`](https://pkg.go.dev/github.com/SherClockHolmes/webpush-go#GenerateVAPIDKeys) from [webpush-go](https://github.com/SherClockHolmes/webpush-go/).

## Usage
```go
package main

import (
    "context"
    "log"

    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/webpush"
)

const vapidPublicKey = "" // Add a vapidPublicKey
const vapidPrivateKey = "" // Add a vapidPrivateKey
const subscription = `` // JSON string of the subscription object

func main() {
    webpushSvc = webpush.New(vapidPublicKey, vapidPrivateKey)

    err := webpushSvc.AddReceivers([]byte(subscription))
    if err != nil {
      log.Fatalf("could not add recivier: %v", err)
    }

    notifier := notify.New()

    notifier.UseServices(webpushSvc)
    if err := notifier.Send(context.Background(), "TEST", "Message using golang notifier library"); err != nil {
        log.Fatalf("notifier.Send() failed: %s", err.Error())
    }

    log.Println("Notification sent")
}
```
