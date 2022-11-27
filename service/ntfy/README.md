# Ntfy

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/ntfy)

## Usage

```go
package main

import (
	"context"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/ntfy"
)

func main() {
	// Create a ntfy service. You can use the
	// `ntfy.NewWithServers` function to create a service with a custom server.
	//ntfyService := ntfy.NewWithServers(ntfy.DefaultServerURL)

	// Or use `ntfy.New` to create a service with the default server.
	ntfyService := ntfy.New()

	// Tell our notifier to use the bark service.
	notify.UseServices(ntfyService)

	content := `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"tags": ["warning","cd"],
		"priority": 4,
		"attach": "https://filesrv.lan/space.jpg",
		"filename": "diskspace.jpg",
		"click": "https://homecamera.lan/xasds1h2xsSsa/",
		"actions": [{ "action": "view", "label": "Admin panel", "url": "https://filesrv.lan/admin" }]
	}`

	// Send a test message.
	err := notify.Send(
		context.Background(),
		"pushkar",
		content,
	)

	if err != nil {
		panic(err)
	}
}
```
