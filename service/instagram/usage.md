# Instagram Usage

Ensure that you have already navigated to your GOPATH and installed the following packages:

* `go get -u github.com/nikoksr/notify`
* `go get github.com/ahmdrz/goinsta/v2` - You might need this one too


## Sample Code

```go
package main

import (
    "github.com/nikoksr/notify"
    "github.com/nikoksr/notify/service/instagram"
)

func main() {

    notifier := notify.New()

    // Assuming you already have an Instagram account
    // Provide your Instagram username and password
    service := instagram.New("Username", "Password")

    // Passing usernames as receivers for our messages.
    service.AddReceivers("friend1", "friend2")

    // Tell our notifier to use the Instagram service. You can repeat the above process
    // for as many services as you like and just tell the notifier to use them.
    notifier.UseServices(service)

    // Send a message
    err := notifier.Send(
        "Hello\n",
        "I am a bot written in Go!",
    )

    if err != nil {
        panic(err)
    }

    // This is Instagram specific, logout after you are done to close the session.
    service.Logout()
}
```