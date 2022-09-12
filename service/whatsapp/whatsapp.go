package whatsapp

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/Rhymen/go-whatsapp"
	"github.com/pkg/errors"
)

const (
	qrLoginWaitTime = 3 * time.Second
	clientTimeout   = 5 * time.Second
)

var sessionFilePath = filepath.Join(os.TempDir(), "whatsappSession.gob")

// whatsappClient abstracts go-whatsapp for writing unit tests
//
//go:generate mockery --name=whatsappClient --output=. --case=underscore --inpackage
type whatsappClient interface {
	Login(qrChan chan<- string) (whatsapp.Session, error)
	RestoreWithSession(session whatsapp.Session) (whatsapp.Session, error)
	Send(msg interface{}) (string, error)
}

// Service encapsulates the WhatsApp client along with internal state for storing contacts.
type Service struct {
	client   whatsappClient
	contacts []string
}

// New returns a new instance of a WhatsApp notification service.
func New() (*Service, error) {
	client, err := whatsapp.NewConn(clientTimeout)
	if err != nil {
		return nil, err
	}

	return &Service{
		client:   client,
		contacts: []string{},
	}, nil
}

// LoginWithSessionCredentials provides helper for authentication using whatsapp.Session credentials.
func (s *Service) LoginWithSessionCredentials(clientID, clientToken, serverToken, wid string, encKey, macKey []byte) error {
	session := whatsapp.Session{
		ClientId:    clientID,
		ClientToken: clientToken,
		ServerToken: serverToken,
		Wid:         wid,
		EncKey:      encKey,
		MacKey:      macKey,
	}

	session, err := s.client.RestoreWithSession(session)
	if err != nil {
		return fmt.Errorf("restoring session failed: %w", err)
	}

	// Save the updated session for future use without login.
	err = writeSession(&session)
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	return nil
}

// LoginWithQRCode provides helper for authentication using QR code on terminal.
// Refer: https://github.com/Rhymen/go-whatsapp#login for more information.
func (s *Service) LoginWithQRCode() error {
	session, err := readSession()
	if err == nil {
		session, err = s.client.RestoreWithSession(session)
		if err != nil {
			return fmt.Errorf("restoring session failed: %w", err)
		}
	} else {
		// No saved session found; need to login again.
		qr := make(chan string)
		go func() {
			terminal := qrcodeTerminal.New()
			terminal.Get(<-qr).Print()
		}()

		session, err = s.client.Login(qr)
		if err != nil {
			return fmt.Errorf("error during login: %w", err)
		}
	}

	err = writeSession(&session)
	if err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	<-time.After(qrLoginWaitTime)

	return nil
}

// readSession helps load saved WhatsApp session from local file system.
func readSession() (whatsapp.Session, error) {
	session := whatsapp.Session{}
	file, err := os.Open(sessionFilePath)
	if err != nil {
		return session, err
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			if err != nil {
				err = errors.Wrap(err, cerr.Error())
			} else {
				err = cerr
			}
		}
	}()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&session)
	if err != nil {
		return session, err
	}

	return session, nil
}

// writeSession helps save WhatsApp session to local file system.
func writeSession(session *whatsapp.Session) error {
	file, err := os.Create(sessionFilePath)
	if err != nil {
		return err
	}
	defer func() {
		cerr := file.Close()
		if cerr != nil {
			if err != nil {
				err = errors.Wrap(err, cerr.Error())
			} else {
				err = cerr
			}
		}
	}()

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
func (s *Service) Send(ctx context.Context, subject, message string) error {
	msg := whatsapp.TextMessage{
		Text: subject + "\n" + message,
	}

	for _, contact := range s.contacts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg.Info = whatsapp.MessageInfo{
				RemoteJid: contact + "@s.whatsapp.net",
			}

			_, err := s.client.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to WhatsApp contact '%s'", contact)
			}
		}
	}

	return nil
}
