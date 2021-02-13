# WhatsApp

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/whatsapp)

## Prerequisites

You will need a registered WhatsApp phone number to be used as source for sending WhatsApp messages.

## Usage

In the current implementation, authentication requires scanning QR code from terminal using a registered WhatsApp device. Please refer [Login (go-whatsapp)](https://github.com/Rhymen/go-whatsapp#login) and [sigalor/whatsapp-web-reveng](https://github.com/sigalor/whatsapp-web-reveng) for more details.

```go
package main

import (
        "log"

        "github.com/nikoksr/notify"
        "github.com/nikoksr/notify/service/whatsapp"
)

func main() {
        whatsappSvc, err := whatsapp.New()
        if err != nil {
                log.Fatalf("whatsapp.New() failed: %s", err.Error())
        }

        whatsappSvc.AddReceivers("Contact1")

        notifier := notify.New()
        notifier.UseServices(whatsappSvc)

        err = notifier.Send("subject", "message")
        if err != nil {
                log.Fatalf("notifier.Send() failed: %s", err.Error())
        }

        log.Println("notification sent")
}
```
