/*
Package viber provides a service for sending messages to viber.

Usage:

	package main

	import (
	    "context"
	    "log"

	    "github.com/nikoksr/notify"
	    "github.com/nikoksr/notify/service/webpush"
	)

	const vapidPublicKey = "" // Add a vapidPublicKey
	const vapidPrivateKey = "" // Add a vapidPrivateKey
  const subscription = `` // JSON string of the subscription object

	func main() {
      webpushSvg = webpush.New(vapidPublicKey, vapidPrivateKey)

      err := webpushSvc.AddReceivers([]byte(subscription))
	    if err != nil {
        log.Fatalf("could not add recivier: %v", err)
	    }

	    notifier := notify.New()

	    notifier.UseServices(webpushSvc)
	    if err := notifier.Send(context.Background(), "TEST", "Message using golang notifier library"); err != nil {
	        log.Fatalf("notifier.Send() failed: %s", err.Error())
	    }

	    log.Println("Notification sent")
	}
*/
package webpush
