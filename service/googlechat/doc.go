/*
Package googlechat provides message notification integration sent to multiple
spaces within a Google Chat Application.

Disclaimer:
In order to integrate `notify` with a Google Chat Application, an "Application Default
Credentials" file must be supplied.

For more information on Google Application credential JSON files see:
https://cloud.google.com/docs/authentication/application-default-credentials

Usage:
    package main

    import (
        "context"
        "fmt"
        "log"
        "strings"

        "github.com/nikoksr/notify"
        "github.com/nikoksr/notify/service/googlechat"
        "google.golang.org/api/chat/v1"
        "google.golang.org/api/option"
    )

    func main() {
        // only basic text messages with subject and message is supported at this time.
        ctx := context.Background()

        withCred := option.WithCredentialsFile("credentials.json")
        withSpacesScope := option.WithScopes("https://www.googleapis.com/auth/chat.spaces") 
        
        listSvc, err := chat.NewService(ctx, withCred, withSpacesScope)
        spaces, err := listSvc.Spaces.List().Do()

        if err != nil {
            log.Fatalf("svc.Spaces.List().Do() failed: %s", err.Error())
        }
        sps := make([]string, 0)
        for _, space := range spaces.Spaces {
        fmt.Printf("space %s\n", space.DisplayName)
            name := strings.Replace(space.Name, "spaces/", "", 1)
            sps = append(sps, name)
        }
        
        msgSvc, err := googlechat.New(withCred)
        if err != nil {
            log.Fatalf("googlechat.New() failed: %s", err.Error())
        }

        msgSvc.AddReceivers(sps...)

        notifier := notify.New()

        notifier.UseServices(msgSvc)

        fmt.Printf("sending message to %d spaces\n", len(sps))
        err = notifier.Send(ctx, "subject", "message")
        if err != nil {
            log.Fatalf("notifier.Send() failed: %s", err.Error())
        }

        log.Println("notification sent")
    }
*/
package googlechat
