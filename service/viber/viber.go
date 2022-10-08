package viber

import (
	"context"

	vb "github.com/mileusna/viber"
	"github.com/pkg/errors"
)

//go:generate mockery --name=viberClient --output=. --case=underscore --inpackage
type viberClient interface {
	SetWebhook(url string, eventTypes []string) (vb.WebhookResp, error)
	SendTextMessage(receiver, msg string) (uint64, error)
}

// Viber struct holds necessary fields to communicate with Viber API
type Viber struct {
	Client            viberClient
	SubscribedUserIDs []string
}

// Compile-time check to ensure that vb.Viber implements the viberClient interface.
var _ viberClient = new(vb.Viber)

// New returns a new instance of Viber notification service
func New(appKey, senderName, senderAvatar string) *Viber {
	return &Viber{
		Client:            vb.New(appKey, senderName, senderAvatar),
		SubscribedUserIDs: []string{},
	}
}

// AddReceivers receives subscribed user IDs then add them to internal receivers list
func (v *Viber) AddReceivers(subscribedUserIDs ...string) {
	v.SubscribedUserIDs = append(v.SubscribedUserIDs, subscribedUserIDs...)
}

func (v *Viber) SetWebhook(webhookURL string) error {
	_, err := v.Client.SetWebhook(webhookURL, []string{})
	return err
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
