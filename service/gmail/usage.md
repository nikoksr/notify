# Gmail Usage

Ensure that you have already navigated to your GOPATH and installed the following packages:

* `go get -u github.com/nikoksr/notify`
* `go get golang.org/x/oauth2`
* `go get google.golang.org/api/gmail/v1`

## Steps for Gmail API

These are general and very high level instructions

1. Create or select a Google Cloud project
2. Enable the Gmail API for that project
3. Configure an OAuth consent screen
4. Create OAuth 2.0 credentials for your application
5. Request a token with a Gmail sending scope such as `https://www.googleapis.com/auth/gmail.send`
6. Use the authenticated email address as the sender address
7. Now you should be good to use the code below

## Sample Code

```go
package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/gmail"
	"golang.org/x/oauth2"
)

func main() {
	notifier := notify.New()

	token := &oauth2.Token{AccessToken: "ACCESS_TOKEN"}
	tokenSource := oauth2.StaticTokenSource(token)

	gmailService, err := gmail.New(tokenSource, "sender@example.com")
	if err != nil {
		log.Fatal(err)
	}

	gmailService.AddReceivers("alice@example.com", "bob@example.com")

	notifier.UseServices(gmailService)

	_ = notifier.Send(
		context.Background(),
		"Hello!",
		"I am a bot written in Go!",
	)
}
```
