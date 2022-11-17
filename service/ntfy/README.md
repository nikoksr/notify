# Ntfy

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/ntfy)

## Usage

```go
package main

import (
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/ntfy"
	"golang.org/x/net/context"
)

func main() {

	notifier := notify.New()
	ntfyService := ntfy.New()
	notifier.UseServices(ntfyService)

	jsonBody := `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"tags": ["warning","cd"],
		"priority": 4,
		"attach": "https://filesrv.lan/space.jpg",
		"filename": "diskspace.jpg",
		"click": "https://homecamera.lan/xasds1h2xsSsa/",
		"actions": [{ "action": "view", "label": "Admin panel", "url": "https://filesrv.lan/admin" }]
	}`

	// Send a message
	err := notifier.Send(
		context.Background(),
		"pushkar",                      // this is topic name
		jsonBody,
	)

	if err != nil {
		panic(err)
	}

}

```
