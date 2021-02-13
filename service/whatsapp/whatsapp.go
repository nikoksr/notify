package whatsapp

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	whatsapp "github.com/Rhymen/go-whatsapp"
	"github.com/pkg/errors"
)

const (
	loginWaitTimeSeconds = 3
	clientTimeout        = 5
)

// whatsappClient abstracts go-whatsapp for writing unit tests
type whatsappClient interface {
	Send(msg interface{}) (string, error)
}

// Service encapsulates the WhatsApp client along with internal state for storing contacts.
type Service struct {
	client   whatsappClient
	contacts []string
}

// New returns a new instance of a WhatsApp notification service.
func New() (*Service, error) {
	client, err := whatsapp.NewConn(clientTimeout * time.Second)
	if err != nil {
		return nil, err
	}

	err = login(client)
	if err != nil {
		return nil, err
	}

	<-time.After(loginWaitTimeSeconds * time.Second)

	s := &Service{
		client:   client,
		contacts: []string{},
	}
	return s, nil
}

// login helps with the WhatsApp authentication process.
// Refer: https://github.com/Rhymen/go-whatsapp#login for more information.
func login(client *whatsapp.Conn) error {
	session, err := readSession()
	if err == nil {
		session, err = client.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring session failed: %v", err)
		}
	} else {
		// No saved session found; need to login again.
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()

		session, err = client.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %v", err)
		}
	}

	err = writeSession(&session)
	if err != nil {
		return fmt.Errorf("error saving session: %v", err)
	}

	return nil
}

// readSession helps load saved WhatsApp session from local file system.
func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return session, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

// writeSession helps save WhatsApp session to local file system.
func writeSession(session *whatsapp.Session) error {
	file, err := os.Create(os.TempDir() + "/whatsappSession.gob")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(*session)
	if err != nil {
		return err
	}

	return nil
}

// AddReceivers takes WhatsApp contacts and adds them to the internal contacts list. The Send method will send
// a given message to all those contacts.
func (s *Service) AddReceivers(contacts ...string) {
	s.contacts = append(s.contacts, contacts...)
}

// Send takes a message subject and a message body and sends them to all previously set contacts.
func (s *Service) Send(subject, message string) error {
	if len(s.contacts) == 0 {
		return fmt.Errorf("no contacts added as receivers")
	}

	msgText := subject + "\n" + message
	for _, c := range s.contacts {
		contact := c + "@s.whatsapp.net"
		msg := whatsapp.TextMessage{
			Info: whatsapp.MessageInfo{
				RemoteJid: contact,
			},
			Text: msgText,
		}

		if _, err := s.client.Send(msg); err != nil {
			return errors.Wrapf(err, "failed to send message to WhatsApp contact '%s'", c)
		}
	}

	return nil
}
