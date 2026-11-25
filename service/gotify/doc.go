/*
Package gotify provides message notification integration for Gotify.

Usage:

	package main

	import (
		"context"
		"log"

		"github.com/nikoksr/notify"
		"github.com/nikoksr/notify/service/gotify"
	)

	func main() {
		// Example 1: Send message with default priority.
		// Create new gotify service with token and service url.
		gotifyService:=New("gotify_token","gotify_service_url")
        
		// Use context.Background() if you want to send a simple notification message.
		ctx := context.Background()
		// Send a simgple message
		gotifyService.Send(context.Background(),"test2","test1234")

		// Example 2: Send message with custom priority.
		// Create new gotify service with token, service url and priority.
		gotifyService:=NewWithPriority("gotify_token", "gotify_service_url", 10)
        
		// Use context.Background() if you want to send a simple notification message.
		ctx := context.Background()
		// Send a simgple message with priority 10
		gotifyService.Send(context.Background(),"test2","test1234")

		log.Println("notification sent")
	}
*/
package gotify
