# Google Chat

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/nikoksr/notify/service/googlechat)

## Prerequisites

In order to integrate `notify` with a Google Chat Application, an "Application Default
Credentials" file must be supplied.

For more information on Google Application credential JSON files see:
https://cloud.google.com/docs/authentication/application-default-credentials

Normally Application Default Credentials as an environment variable would be
supported. However, to remain consistent with other `notify` services, the path to
a valid credential configuration JSON file or service account key JSON file must be
passed in as a parameter of `New` (See Usage for an example).

a example service account key JSON file has been provided in this directory
`example_credentials.json` which takes the following shape:

```json
{
  "type": "service_account",
  "project_id": "",
  "private_key_id": "",
  "private_key": "",
  "client_email": "",
  "client_id": "",
  "auth_uri": "",
  "token_uri": "",
  "auth_provider_x509_cert_url": "",
  "client_x509_cert_url": ""
}
```

## Usage:

```go
package main

import (
    "context"
    "log"

    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/googlechat"
)

func main() {
    gChatSvc, err := googlechat.New("path/to/config_file.json")
    if err != nil {
        log.Fatalf("googlechat.New() failed: %s", err.Error())
    }

    gChatSvc.AddReceivers("office_space")

    notifier := notify.New()
    notifier.UseServices(gChatSvc)

    // only basic text messages with subject and message is supported at
    // this time.
    ctx := context.Background()

    err = notifier.Send(ctx, "subject", "message")
    if err != nil {
        log.Fatalf("notifier.Send() failed: %s", err.Error())
    }

    log.Println("notification sent")
}
```
