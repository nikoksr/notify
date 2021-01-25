package notify

import (
	"github.com/pkg/errors"
)

const defaultDisabled = false // Notifier is enabled by default

// Notify is the central struct for managing notification services and sending messages to them.
type Notify struct {
	Disabled  bool
	notifiers []Notifier
}

// ErrSendNotification signals that the notifier failed to send a notification.
var ErrSendNotification = errors.New("Send notification")

// Notifier defines the behavior for notification services. The Send command simply sends a message string to the
// internal destination Notifier. E.g for telegram it sends the message to the specified group chat.
type Notifier interface {
	Send(string, string) error
}

// New returns a new instance of Notify. Defaulting to being not disabled and using the pseudo notification
// service under the hood.
func New() *Notify {
	notifier := &Notify{
		Disabled: defaultDisabled,
	}

	// Use the pseudo Notifier to prevent from nil reference bugs when using the Notify Notifier. In case no notifiers
	// are provided or the creation of all other notifiers failed, the pseudo Notifier will be used under the hood
	// doing nothing but preventing nil-reference errors.
	notifier.usePseudo()

	return notifier
}
