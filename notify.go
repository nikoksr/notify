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
var ErrSendNotification = errors.New("send notification")

// New returns a new instance of Notify. Defaulting to being not disabled.
func New() *Notify {
	notifier := &Notify{
		Disabled:  defaultDisabled,
		notifiers: []Notifier{},
	}

	return notifier
}
