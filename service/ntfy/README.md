# Ntfy

## Prerequisites

To use the Ntfy service, you need:

1. A running ntfy server (self-hosted or use ntfy.sh)
2. Topic names to send notifications to

## Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/ntfy"
)

func main() {
    // Create ntfy service (using public ntfy.sh server)
    ntfyService := ntfy.New("ntfy.sh")
    
    // Add topics to send notifications to
    ntfyService.AddReceivers("my-app-alerts", "server-monitoring")
    
    // Add service to notify
    notify.UseServices(ntfyService)
    
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

## Self-hosted Server

```go
// Use your own ntfy server
ntfyService := ntfy.New("https://ntfy.example.com")
// or
ntfyService := ntfy.New("http://localhost:8080")
```

## Features

- ✅ Send to multiple topics
- ✅ Custom server support
- ✅ HTTP/HTTPS support
- ✅ Context cancellation
- ✅ Proper error handling

## Links

- [Ntfy Documentation](https://docs.ntfy.sh/)
- [Ntfy GitHub](https://github.com/binwiederhier/ntfy)
- [Public Server](https://ntfy.sh/)