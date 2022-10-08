package viber

import (
	"context"

	vb "github.com/mileusna/viber"
	"github.com/pkg/errors"
)

// Viber struct holds necessary fields to communicate with Viber API
type Viber struct {
	Client            *vb.Viber
	SubscribedUserIDs []string
}

// ViberOption struct holds an optional function to create Viber instance
type ViberOption func(v *Viber)

// New returns a new instance of Viber notification service
func New(appKey, webhookURL, senderName string, opts ...ViberOption) (*Viber, error) {
	viber := &Viber{
		Client:            vb.New(appKey, senderName, ""),
		SubscribedUserIDs: []string{},
	}

	for _, opt := range opts {
		opt(viber)
	}

	_, err := viber.Client.SetWebhook(webhookURL, []string{})
	if err != nil {
		return nil, err
	}

	return viber, nil
}

// WithSenderAvatar function to add senderAvatar to the Viber client
func WithSenderAvatar(senderAvatar string) ViberOption {
	return func(v *Viber) {
		v.Client.Sender.Avatar = senderAvatar
	}
}

// AddReceivers receives subscribed user IDs then add them to internal receivers list
func (v *Viber) AddReceivers(subscribedUserIDs ...string) {
	v.SubscribedUserIDs = append(v.SubscribedUserIDs, subscribedUserIDs...)
}

// Send takes a message subject and a message body and sends them to all previously set userIds
func (v *Viber) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, subscribedUserID := range v.SubscribedUserIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, err := v.Client.SendTextMessage(subscribedUserID, fullMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to User ID '%s'", subscribedUserID)
			}
		}
	}

	return nil
}
