package pushbullet

import (
	"context"

	"github.com/cschomburg/go-pushbullet"
	"github.com/pkg/errors"
)

// SMS struct holds necessary data to communicate with the Pushbullet SMS API.
type SMS struct {
	client           *pushbullet.Client
	deviceIdentifier string
	phoneNumbers     []string
}

// NewSMS returns a new instance of a SMS notification service
// tied to an SMS capable device. deviceNickname is the
// Pushbullet nickname of the sms capable device from which messages are sent.
// (https://help.pushbullet.com/articles/how-do-i-send-text-messages-from-my-computer/).
// For more information about Pushbullet api token:
//
//	-> https://docs.pushbullet.com/#api-overview
func NewSMS(apiToken, deviceNickname string) (*SMS, error) {
	client := pushbullet.New(apiToken)

	dev, err := client.Device(deviceNickname)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to locate Pushbullet device with nickname '%s'", deviceNickname)
	}

	sms := &SMS{
		client:           client,
		deviceIdentifier: dev.Iden,
		phoneNumbers:     []string{},
	}

	return sms, nil
}

// AddReceivers takes phone numbers and adds them to the internal phoneNumbers list. The Send method will send
// a given message to all registered phone numbers.
func (sms *SMS) AddReceivers(phoneNumbers ...string) {
	sms.phoneNumbers = append(sms.phoneNumbers, phoneNumbers...)
}

// Send takes a message subject and a message body and sends them to all phone numbers.
// see https://help.pushbullet.com/articles/how-do-i-send-text-messages-from-my-computer/
func (sms SMS) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title
	user, err := sms.client.Me()
	if err != nil {
		return errors.Wrapf(err, "failed to find valid pushbullet user")
	}

	for _, phoneNumber := range sms.phoneNumbers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err = sms.client.PushSMS(user.Iden, sms.deviceIdentifier, phoneNumber, fullMessage)
			if err != nil {
				return errors.Wrapf(err, "failed to send SMS message to %s via Pushbullet", phoneNumber)
			}
		}
	}

	return nil
}
