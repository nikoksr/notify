/*
Package homeassistant provides message notification integration for Home Assistant (https://www.home-assistant.io).

Home Assistant is an open source home automation platform that puts local control and privacy first.
This service allows sending notifications through Home Assistant's notification system.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/homeassistant"
	)

	func main() {
		haService := homeassistant.New("http://homeassistant.local:8123", "your-long-lived-access-token")

		notify.UseServices(haService)

		err := notify.Send(context.Background(), "Subject", "Message")
		if err != nil {
			log.Fatal(err)
		}
	}
*/
package homeassistant
