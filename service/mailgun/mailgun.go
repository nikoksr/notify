package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/pkg/errors"
)

// Mailgun struct holds necessary data to communicate with the Mailgun API.
type Mailgun struct {
	client            mailgun.Mailgun
	senderAddress     string
	receiverAddresses []string
}

// New returns a new instance of a Mailgun notification service.
// You will need a Mailgun API key and domain name.
// See https://documentation.mailgun.com/en/latest/
func New(domain, apiKey, senderAddress string, opts ...Option) *Mailgun {
	m := &Mailgun{
		client:            mailgun.NewMailgun(domain, apiKey),
		senderAddress:     senderAddress,
		receiverAddresses: []string{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *Mailgun) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (m Mailgun) Send(ctx context.Context, subject, message string) error {
	mailMessage := m.client.NewMessage(m.senderAddress, subject, message, m.receiverAddresses...)

	_, _, err := m.client.Send(ctx, mailMessage)
	if err != nil {
		return errors.Wrap(err, "failed to send mail using Mailgun service")
	}

	return nil
}
