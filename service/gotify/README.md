# Gotify

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/gotify)

## Prerequisites

Navigate to your Gotify Service, login with your account, and create a new application.
You will find the `Token` for the new created application.

If you want to find more details, please go to [gotify documention](https://gotify.net/docs/index).

## Usage

### Send message with default priority(5)
```go
package main

import (
  "context"
  "log"

  "github.com/nikoksr/notify"
  "github.com/nikoksr/notify/service/gotify"
)

func main() {
  // Create new gotify service with token and service url.
  gotifyService:=gotify.New("gotify_token","gotify_service_url")

  // Use context.Background() if you want to send a simple notification message.
  ctx := context.Background()
  // Send a simgple message
  gotifyService.Send(ctx,"test2","test1234")

  log.Println("notification sent")
}
```

### Send message with custom priority
```go
package main

import (
  "context"
  "log"

  "github.com/nikoksr/notify"
  "github.com/nikoksr/notify/service/gotify"
)

func main() {
  // Create new gotify service with token and service url.
  gotifyService:=gotify.NewWithPriority("gotify_token", "gotify_service_url", 10)

  // Use context.Background() if you want to send a simple notification message.
  ctx := context.Background()
  // Send a simgple message
  gotifyService.Send(ctx,"test2","test1234")

  log.Println("notification sent")
}
```
