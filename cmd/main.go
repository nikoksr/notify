package main

import (
	"context"
	"log"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/viber"
)

const appKey = "4fedd3b160e7dc77-8378815e812898c4-f8bb7ce536c7b68f"
const webhookURL = "https://script.google.com/macros/s/AKfycbzMMYywzyZ5eFDo6rNowWPAmo1BMTbQ3QgT8Mx-oVun0YhnE-L4AJ5TkgLZMbcZDIPa/exec"
const senderName = "vibersofyana"

func main() {
	viberSvc, err := viber.New(appKey, webhookURL, senderName)
	if err != nil {
		log.Fatal(err)
	}

	viberSvc.AddReceivers("Ocounr/S1Ub8FOS3PQdNyw==", "Ocounr/S1Ub8FOS3PQdNyw==")

	notifier := notify.New()
	notifier.UseServices(viberSvc)

	if err := notifier.Send(context.Background(), "TEST", "Message using notifier"); err != nil {
		log.Fatalf("notifier.Send() failed: %s", err.Error())
	}

	log.Println("Notification sent")
}
