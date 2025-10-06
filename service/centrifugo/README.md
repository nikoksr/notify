# Centrifugo Notification Service

This service allows you to send real-time notifications to [Centrifugo](https://centrifugal.dev/) channels using the [centrifuge-go](https://github.com/centrifugal/centrifuge-go) client.

## Features
- Send messages to Centrifugo channels over WebSocket
- Supports subject and message (concatenated)
- Easy integration with the notify library

## Usage

```go
import (
	"context"
	centrifugo "github.com/nikoksr/notify/service/centrifugo"
)

// Create a new Centrifugo service
svc, err := centrifugo.New("ws://localhost:8000/connection/websocket", "your-channel", "your-jwt-token-if-any")
if err != nil {
	panic(err)
}
defer svc.Close()

// Send a notification
err = svc.Send(context.Background(), "Subject", "Hello from notify!")
```

## Links
- [Centrifugo Documentation](https://centrifugal.dev/docs/)
- [Go Client: centrifuge-go](https://github.com/centrifugal/centrifuge-go)
