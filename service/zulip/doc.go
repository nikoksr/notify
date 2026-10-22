/*
Package zulip provides message integration for the zulip service.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/zulip"
	)

	func main() {
	    zulipSvc, err := zulip.New("server-base-url", "bot-email", "api-key")
		if err != nil {
			log.Fatalf("zulip.New() failed: %s", err.Error())
		}

	    zulipSvc.AddReceivers(Direct("test2@gmail.com"), Stream("stream", "topic"))

		notifier := notify.New()
		notifier.UseServices(zulipSvc)

		err = notifier.Send(context.Background(), "subject", "message")
		if err != nil {
			log.Fatalf("notifier.Send() failed: %s", err.Error())
		}

		log.Println("notification sent")
	}
*/
package zulip
