package mail

import (
	"context"
	"crypto/tls"
	"net"
	"net/textproto"
	"strconv"
	"time"

	"gopkg.in/mail.v2"
)

// Mail struct holds necessary data to send emails.
type Mail struct {
	usePlainText      bool
	senderAddress     string
	smtpHostAddr      string
	user              string
	pass              string
	receiverAddresses []string
	headers           textproto.MIMEHeader
	tls               *tls.Config
}

// New returns a new instance of a Mail notification service.
func New(senderAddress, smtpHostAddress string) *Mail {
	return &Mail{
		usePlainText:      false,
		senderAddress:     senderAddress,
		smtpHostAddr:      smtpHostAddress,
		receiverAddresses: []string{},
		headers:           textproto.MIMEHeader{},
	}
}

// BodyType is used to specify the format of the body.
type BodyType int

const (
	// PlainText is used to specify that the body is plain text.
	PlainText BodyType = iota
	// HTML is used to specify that the body is HTML.
	HTML
)

// AddReceivers takes email addresses and adds them to the internal address list. The Send method will send
// a given message to all those addresses.
func (m *Mail) AddReceivers(addresses ...string) {
	m.receiverAddresses = append(m.receiverAddresses, addresses...)
}

func (m *Mail) AddAuthentication(user string, pass string) {
	m.user = user
	m.pass = pass
}

// BodyFormat can be used to specify the format of the body.
// Default BodyType is HTML.
func (m *Mail) BodyFormat(format BodyType) {
	switch format {
	case PlainText:
		m.usePlainText = true
	default:
		m.usePlainText = false
	}
}

func (m *Mail) AddHeader(name, value string) {
	if m.headers == nil {
		m.headers = textproto.MIMEHeader{}
	}
	m.headers.Add(name, value)
}

func (m *Mail) InsecureSkipVerify(enable bool) {
	if m.tls == nil {
		host, _, _ := net.SplitHostPort(m.smtpHostAddr)
		m.tls = &tls.Config{ServerName: host}
	}
	m.tls.InsecureSkipVerify = enable
}

func (m *Mail) newEmail(subject, message string) *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("From", m.senderAddress)
	msg.SetHeader("To", m.receiverAddresses...)
	msg.SetHeader("Subject", subject)
	if m.usePlainText {
		msg.SetBody("text/plain", message)
	} else {
		msg.SetBody("text/html", message)

	}

	for name, values := range m.headers {
		for _, value := range values {
			msg.SetHeader(name, value)
		}
	}

	return msg
}

// Send takes a message subject and a message body and sends them to all previously set chats. Message body supports
// html as markup language.
func (m Mail) Send(ctx context.Context, subject, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	msg := m.newEmail(subject, message)

	host, portStr, err := net.SplitHostPort(m.smtpHostAddr)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	d := mail.NewDialer(host, port, m.user, m.pass)
	if m.tls != nil {
		d.TLSConfig = m.tls
	}
	if deadline, ok := ctx.Deadline(); ok {
		timeout := time.Until(deadline)
		if timeout > 0 {
			d.Timeout = timeout
		} else {
			return context.DeadlineExceeded
		}
	}

	err = d.DialAndSend(msg)
	return err
}
