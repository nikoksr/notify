# Zulip

## Prerequisites

To use the Zulip service, you need:

1. A Zulip server (cloud or self-hosted)
2. A bot account with API credentials
3. Email and API key for the bot

## Usage

### Basic Usage
```go
package main

import (
    "context"
    "log"
    
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/zulip"
)

func main() {
    // Create Zulip service
    zulipService := zulip.New("https://company.zulipchat.com", "bot@company.com", "your-api-key")
    
    // Add service to notify
    notify.UseServices(zulipService)
    
    // Send notification to default stream
    err := notify.Send(
        context.Background(),
        "Build Status",
        "Deployment completed successfully",
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Stream Messages
```go
// Send to specific stream with topic
err := zulipService.SendToStream(ctx, "development", "CI/CD", "Build #123 passed")
```

### Direct Messages
```go
// Send private message to users
err := zulipService.SendDirectMessage(ctx, "alice@company.com,bob@company.com", "Meeting reminder")
```

## Getting API Credentials

1. Login to your Zulip instance
2. Go to Settings → Account & privacy
3. Scroll to "API key" section
4. Click "Show/change your API key"
5. Copy your email and API key

### Creating a Bot Account
1. Go to Settings → Your bots
2. Click "Add a new bot"
3. Choose bot type and fill details
4. Copy the bot email and API key

## Configuration Examples

### Zulip Cloud
```go
zulipService := zulip.New("https://company.zulipchat.com", "bot@company.com", "api-key")
```

### Self-hosted Zulip
```go
zulipService := zulip.New("https://chat.company.com", "bot@company.com", "api-key")
```

## Features

- ✅ Stream messages with topics
- ✅ Direct/private messages
- ✅ Multiple recipient support
- ✅ HTTP/HTTPS support
- ✅ Basic authentication
- ✅ Context cancellation
- ✅ Comprehensive error handling

## Links

- [Zulip Documentation](https://zulip.com/help/)
- [Zulip API Documentation](https://zulip.com/api/)
- [Bot Setup Guide](https://zulip.com/help/add-a-bot-or-integration)