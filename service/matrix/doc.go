/*
Package matrix provides message notification integration for Matrix.

Usage:

	package main

import (

	"context"
	"log"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/matrix"

)

		func main() {
		  matrixSvc, err := matrix.New("user-id", "room-id", "home-server", "access-token")
		  if err != nil {
		    log.Fatalf("matrix.New() failed: %s", err.Error())
		  }

		  notifier := notify.New()
		  notifier.UseServices(matrixSvc)

		  err = notifier.Send(context.Background(), "", "message")
		  if err != nil {
		    log.Fatalf("notifier.Send() failed: %s", err.Error())
		  }

		  log.Println("notification sent")
	    }
*/
package matrix
