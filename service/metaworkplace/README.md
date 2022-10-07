## Meta Workplace Service

### Usage
* You must have a Meta Workplace access token to use this.
* You must be part of an existing thread/chat with a user to use this.
* If you are not in a thread/chat with a user already you must use AddUsers().
* Unless AddUsers() is called again subsequent calls to SendMessage() will send to the same thread/chat.

```go
package main

import (
	"context"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/metaworkplace"
	"log"
)

func main() {
	// Create a new notifier
	notifier := notify.New()

	accesstoken := "your_access_token"

	// Create a new Workplace service.
	wps := metaworkplace.New(accesstoken)

	wps.AddUsers("some_user_id")
	wps.AddThreads("t_somethreadid", "t_somethreadid")

	notifier.UseServices(wps)

	err := notifier.Send(context.Background(), "", "our first message")
	if err != nil {
		log.Fatal(err)
	}

	err = notifier.Send(context.Background(), "", "our second message")
	if err != nil {
		log.Fatal(err)
	}

	wps.AddUsers("some_other_user_id")

	err = notifier.Send(context.Background(), "", "our third message")
	if err != nil {
		log.Fatal(err)
	}
}
```


#### Please consider the following:
* Meta Workplace service does not allow sending messages to group chats
* Meta Workplace service does not allow sending messages to group posts
* Meta Workplace service does not allow sending messages to users that are not already in a thread/chat with you
* Testing still needs to be added; the scenarios below will all fail. Those without a user/thread id should fail. Those
with a valid user/thread id should succeed and return an error for the invalid user/thread id(s). Here are the scenarios:
  * wps.AddUsers("", "123456789012345")
    wps.AddUsers("", "123456789012345", "")
    wps.AddUsers("")
  * wps.AddThreads("", "123456789012345")
    wps.AddThreads("", "123456789012345", "")
    wps.AddThreads("")

## Contributors
- [Melvin Hillsman](github.com/mrhillsman)
- [Ainsley Clark](github.com/ainsleyclark)
