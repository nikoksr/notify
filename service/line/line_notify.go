package line

import (
	"context"

	"github.com/pkg/errors"
	"github.com/utahta/go-linenotify"
)

// Line Notify struct holds info about client and destination token for communicating with line API
type Notify struct {
	client         *linenotify.Client
	receiverTokens []string
}

// New creates a new instance of Line notify service
// For more info about line notify api:
// -> https://notify-bot.line.me/doc/en/
func NewNotify() *Notify {
	c := linenotify.NewClient()
	l := &Notify{
		client: c,
	}
	return l
}

// AddReceivers receives token then add them to internal receivers list
func (ln *Notify) AddReceivers(receiverTokens ...string) {
	ln.receiverTokens = append(ln.receiverTokens, receiverTokens...)
}

// Send receives message subject and body then sends it to all receivers set previously
// Subject will be on the first line followed by message on the next line
func (ln *Notify) Send(ctx context.Context, subject, message string) error {
	lineMessage := subject + "\n" + message

	for _, receiverToken := range ln.receiverTokens {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := ln.client.NotifyMessage(ctx, receiverToken, lineMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to LINE contact '%s'", receiverToken)
			}
		}
	}

	return nil
}
