package sendgrid

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGrid struct holds necessary data to communicate with the SendGrid API.
type SendGrid struct {
	client            *sendgrid.Client
	senderAddress     string
	senderName        string
	receiverAddresses []string
}

// New returns a new instance of a SendGrid notification service.
// You will need a SendGrid API key.
// See https://sendgrid.com/docs/for-developers/sending-email/api-getting-started/
func New(apiKey, senderAddress, senderName string) *SendGrid {
	client := sendgrid.NewSendClient(apiKey)

	return &SendGrid{
		client:            client,
		senderAddress:     senderAddress,
		senderName:        senderName,
		receiverAddresses: []string{},
	}
}

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (s *SendGrid) AddReceivers(addresses ...string) {
	s.receiverAddresses = append(s.receiverAddresses, addresses...)
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (s SendGrid) Send(subject, message string) error {
	from := mail.NewEmail(s.senderName, s.senderAddress)
	c := mail.NewContent("text/html", message)

	// Create a new personalization instance to be able to add multiple receiver addresses.
	p := mail.NewPersonalization()
	p.Subject = subject

	for _, receiverAddress := range s.receiverAddresses {
		receiverEmail := mail.NewEmail(receiverAddress, receiverAddress)
		p.AddTos(receiverEmail)
	}

	m := mail.NewV3Mail()
	m.AddPersonalizations(p)
	m.AddContent(c)
	m.SetFrom(from)

	resp, err := s.client.Send(m)
	if err != nil {
		return errors.Wrap(err, "failed to send mail using SendGrid service")
	}

	if resp.StatusCode != http.StatusAccepted {
		return errors.New("failed to send mail using SendGrid service")
	}

	return nil
}
