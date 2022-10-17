# Mattermost Usage

Ensure that you have already navigated to your GOPATH and installed the following packages:

* `go get -u github.com/nikoksr/notify`
* `go get github.com/mattermost/mattermost-server/v5` - You might need this one too

## Steps for Mattermost Server

These are general and very high level instructions

1. Create a new Mattermost server / Join existing Mattermost server
2. Make sure your Username/loginID have the OAuth permission scope(s): `create_post`
3. Copy the *Channel ID* of the channel you want to post a message to. You can grab the *Channel ID* in channel info. example: *yfgstwuisnshydhd*
4. Now you should be good to use the code below

## Sample Code

```go
package main

import (
    "os"

    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/mattermost"
)

func main() {

    // Init notifier
    notifier := notify.New()

    // Provide your Mattermost server url
    mattermostService := mattermost.New("https://myserver.cloud.mattermost.com")

    // Provide username as loginID and password to login into above server.
    err := mattermostService.LoginWithCredentials("someone@gmail.com", "somepassword")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Passing a Mattermost channel/chat id as receiver for our messages.
    // Where to send our messages.
    mattermostService.AddReceivers("CHANNEL_ID")

    // Tell our notifier to use the Mattermost service. You can repeat the above process
    // for as many services as you like and just tell the notifier to use them.
    notifier.UseServices(mattermostService)

    // Send a message
    err = notifier.Send(
        context.Background(),
        "Hello :wave:\n",
        "Message written in Go!",
    )
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

}
```
