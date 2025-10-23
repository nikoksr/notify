# Gotify

## Prerequisites

To use the Gotify service, you need:

1. A running Gotify server (self-hosted)
2. An application token from your Gotify server

## Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/gotify"
)

func main() {
    // Create gotify service
    gotifyService := gotify.New("https://gotify.example.com", "your-app-token")
    
    // Add service to notify
    notify.UseServices(gotifyService)
    
    // Send notification
    err := notify.Send(
        context.Background(),
        "Server Alert",
        "High CPU usage detected on server-01",
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Getting App Token

1. Login to your Gotify web interface
2. Go to "Apps" section
3. Create a new application
4. Copy the generated token

## Features

- ✅ Simple HTTP API integration
- ✅ Custom server support
- ✅ HTTP/HTTPS support
- ✅ Context cancellation
- ✅ Priority support (0-10)
- ✅ Proper error handling

## Links

- [Gotify Documentation](https://gotify.net/docs/)
- [Gotify GitHub](https://github.com/gotify/server)
- [API Documentation](https://gotify.net/api-docs)