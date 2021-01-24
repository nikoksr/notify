package notify

import (
	"github.com/pkg/errors"
)

type Notifier struct {
	Disabled bool
	services []Service
}

const defaultDisabled = false

var ErrSendNotification = errors.New("Send notification")

// Service implements a Listen and a Send method. The Listen method makes the notification Service listens
// for external commands and will answer to these commands if supported. For example, our telegram Notifier listens for
// commands like /info and will answer with basic information about the server. The Send command simply sends a
// message string to the internal destination Service. E.g for telegram it sends the message to the specified group
// chat.
type Service interface {
	Send(string, string) error
}

func New() *Notifier {
	notifier := &Notifier{
		Disabled: defaultDisabled,
	}

	// Use the pseudo Service to prevent from nil reference bugs when using the Notifier Service. In case no services
	// are provided or the creation of all other services failed, the pseudo Service will be used under the hood
	// doing nothing but preventing nil-reference errors.
	notifier.usePseudo()

	return notifier
}
