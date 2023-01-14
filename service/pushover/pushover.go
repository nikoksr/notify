// Package pushover implements a Pushover notifier, allowing messages to be sent to multiple recipients and supports
// both users and groups.
package pushover

import (
	"context"

	"github.com/gregdel/pushover"
	"github.com/pkg/errors"
)

//go:generate mockery --name=pushoverClient --output=. --case=underscore --inpackage
type pushoverClient interface {
	SendMessage(*pushover.Message, *pushover.Recipient) (*pushover.Response, error)
}

// Compile-time check to ensure that pushover.Pushover implements the pushoverClient interface.
var _ pushoverClient = new(pushover.Pushover)

// Pushover struct holds necessary data to communicate with the Pushover API.
type Pushover struct {
	client     pushoverClient
	recipients []pushover.Recipient
}

// New returns a new instance of a Pushover notification service.
// For more information about Pushover app token:
//
//	-> https://support.pushover.net/i175-how-do-i-get-an-api-or-application-token
func New(appToken string) *Pushover {
	client := pushover.New(appToken)

	s := &Pushover{
		client:     client,
		recipients: []pushover.Recipient{},
	}

	return s
}

// AddReceivers takes Pushover user/group IDs and adds them to the internal recipient list. The Send method will send
// a given message to all of those recipients.
func (p *Pushover) AddReceivers(recipientIDs ...string) {
	for _, recipient := range recipientIDs {
		p.recipients = append(p.recipients, *pushover.NewRecipient(recipient))
	}
}

// Send takes a message subject and a message body and sends them to all previously set recipients.
func (p Pushover) Send(ctx context.Context, subject, message string) error {
	for i := range p.recipients {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := p.client.SendMessage(
				pushover.NewMessageWithTitle(message, subject),
				&p.recipients[i],
			)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Pushover recipient '%s'", p.recipients[i])
			}
		}
	}
	return nil
}
