/*
Package whatsapp provides message notification integration for WhatsApp.

Usage:

	package main

	import (
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/whatsapp"
	)

	func main() {
		whatsappSvc, err := whatsapp.New()
		if err != nil {
			log.Fatalf("whatsapp.New() failed: %s", err.Error())
		}

		err = whatsappSvc.LoginWithQRCode()
		if err != nil {
			log.Fatalf("whatsappSvc.LoginWithQRCode() failed: %s", err.Error())
		}

		whatsappSvc.AddReceivers("Contact1")

		notifier := notify.New()
		notifier.UseServices(whatsappSvc)

		err = notifier.Send(context.Background(), "subject", "message")
		if err != nil {
			log.Fatalf("notifier.Send() failed: %s", err.Error())
		}

		log.Println("notification sent")
	}
*/
package whatsapp
