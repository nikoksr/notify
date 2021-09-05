package twilio

import (
	"context"
	"net/url"

	"github.com/kevinburke/twilio-go"
	"github.com/pkg/errors"
)

// twilioMessageClient abstracts twilio-go for writing unit tests
type twilioMessageClient interface {
	SendMessage(from, to, body string, mediaURLs []*url.URL) (*twilio.Message, error)
}

// Service encapsulates the Twilio client along with internal state for storing contacts.
type Service struct {
	messageClient twilioMessageClient
	phone         string
	contacts      []string
}

// New returns a new instance of a Twilio message service.
func New(phoneNo, sid, token string) (*Service, error) {
	client := twilio.NewClient(sid, token, nil)

	s := &Service{
		messageClient: client.Messages,
		contacts:      []string{},
		phone:         phoneNo,
	}
	return s, nil
}

// AddReceivers takes Twilio contacts and adds them to the internal contacts list. The Send method will send
// a given message to all those contacts.
func (s *Service) AddReceivers(contacts ...string) {
	s.contacts = append(s.contacts, contacts...)
}

// Send takes a messageClient subject and a messageClient body and sends them to all previously set contacts.
func (s *Service) Send(ctx context.Context, subject, message string) error {

	msg := subject + "\n" + message

	for _, contact := range s.contacts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:

			_, err := s.messageClient.SendMessage(s.phone, contact, msg, nil)
			if err != nil {
				return errors.Wrapf(err, "failed to send messageClient to Twilio contact '%s'", contact)
			}
		}
	}

	return nil
}
