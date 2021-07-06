package twilio

import (
    "context"
    "github.com/kevinburke/twilio-go"
    "github.com/pkg/errors"
)

// Service encapsulates the WhatsApp client along with internal state for storing contacts.
type Service struct {
    client   *twilio.Client
    phone    string
    contacts []string
}

// New returns a new instance of a WhatsApp notification service.
func New(phoneNo, sid, token string) (*Service, error) {
    client := twilio.NewClient(sid, token, nil)

    s := &Service{
        client:   client,
        contacts: []string{},
        phone: phoneNo,
    }
    return s, nil
}

// AddReceivers takes Twilio contacts and adds them to the internal contacts list. The Send method will send
// a given message to all those contacts.
func (s *Service) AddReceivers(contacts ...string) {
    s.contacts = append(s.contacts, contacts...)
}

// Send takes a message subject and a message body and sends them to all previously set contacts.
func (s *Service) Send(ctx context.Context, subject, message string) error {

    msg := subject + "\n" + message

    for _, contact := range s.contacts {
        _, err := s.client.Messages.SendMessage(s.phone, contact, msg, nil)
        if err != nil {
            return errors.Wrapf(err, "failed to send message to Twilio contact '%s'", contact)
        }
    }

    return nil
}
