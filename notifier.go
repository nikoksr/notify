package notify

import "context"

// Notifier defines the behavior for notification services.
//
// The Send function simply sends a subject and a message string to the internal destination Notifier.
//
//	E.g. for telegram.Telegram it sends the message to the specified group chat.
type Notifier interface {
	Send(context.Context, string, string) error
}
