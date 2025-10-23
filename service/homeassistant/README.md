# Home Assistant

## Prerequisites

To use the Home Assistant service, you need:

1. A running Home Assistant instance
2. A long-lived access token from Home Assistant

## Usage

```go
package main

import (
    "context"
    "log"
    
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/homeassistant"
)

func main() {
    // Create Home Assistant service
    haService := homeassistant.New("http://homeassistant.local:8123", "your-long-lived-access-token")
    
    // Add service to notify
    notify.UseServices(haService)
    
    // Send notification
    err := notify.Send(
        context.Background(),
        "Security Alert",
        "Motion detected in living room",
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## Getting Access Token

1. Login to your Home Assistant web interface
2. Go to your Profile (click on your username)
3. Scroll down to "Long-Lived Access Tokens"
4. Click "Create Token"
5. Give it a name and copy the generated token

## Configuration Examples

### Local Instance
```go
haService := homeassistant.New("http://192.168.1.100:8123", "token")
```

### Remote Instance with SSL
```go
haService := homeassistant.New("https://my-ha.duckdns.org", "token")
```

### Nabu Casa Cloud
```go
haService := homeassistant.New("https://xxxxx.ui.nabu.casa", "token")
```

## Features

- ✅ Persistent notifications
- ✅ HTTP/HTTPS support
- ✅ Long-lived access token authentication
- ✅ Context cancellation
- ✅ Proper error handling
- ✅ Local and remote instance support

## Links

- [Home Assistant Documentation](https://www.home-assistant.io/docs/)
- [Home Assistant API](https://developers.home-assistant.io/docs/api/rest/)
- [Notification Services](https://www.home-assistant.io/integrations/#notifications)