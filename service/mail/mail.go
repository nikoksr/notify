package mail

import (
	"github.com/pkg/errors"
	gomail "gopkg.in/mail.v2"
)

const (
	headerFrom      = "From"
	headerTo        = "To"
	headerSubject   = "Subject"
	bodyContentType = "text/plain"
)

type Mail struct {
	client            *gomail.Dialer
	receiverAddresses []string
}

func New(host, userName, password string, port int) (*Mail, error) {
	client := gomail.NewDialer(host, port, userName, password)

	m := &Mail{
		client:            client,
		receiverAddresses: []string{},
	}

	return m, nil
}

func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

func (m Mail) Send(subject, message string) error {
	msg := gomail.NewMessage()

	// Set E-Mail sender
	msg.SetHeader(headerFrom, m.client.Host)

	// Set E-Mail receivers
	msg.SetHeader(headerTo, m.receiverAddresses...)

	// Set E-Mail subject
	msg.SetHeader(headerSubject, subject)

	// Set E-Mail body
	msg.SetBody(bodyContentType, message)

	err := m.client.DialAndSend(msg)
	if err != nil {
		err = errors.Wrap(err, "failed to dial and send mail")
	}

	return err
}
