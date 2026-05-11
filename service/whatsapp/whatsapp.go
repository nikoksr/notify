package whatsapp

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "modernc.org/sqlite" // SQLite driver required by whatsmeow store; imported for side-effects only.
)

// ErrMissingRecipient is returned when Send is called with no recipients configured.
var ErrMissingRecipient = errors.New("whatsapp: missing recipient JID")

// Service implements notify.Notifier for WhatsApp using whatsmeow.
type Service struct {
	client     *whatsmeow.Client
	recipients []types.JID
}

// New returns a new instance of a WhatsApp notification service.
func New() *Service {
	return &Service{}
}

// AddReceivers takes WhatsApp JID strings and appends them to the internal recipient list.
// The expected format is: 6281234567890@s.whatsapp.net.
func (s *Service) AddReceivers(receivers ...string) {
	for _, r := range receivers {
		jid, err := types.ParseJID(r)
		if err != nil {
			continue
		}

		s.recipients = append(s.recipients, jid)
	}
}

func (s *Service) initClient(ctx context.Context, dbPath string) error {
	container, err := sqlstore.New(ctx, "sqlite", "file:"+dbPath+"?_foreign_keys=on", waLog.Noop)
	if err != nil {
		return fmt.Errorf("whatsapp: failed to init store: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return fmt.Errorf("whatsapp: failed to get device: %w", err)
	}

	s.client = whatsmeow.NewClient(deviceStore, waLog.Noop)

	return nil
}

// LoginWithQRCode authenticates via QR code printed to the terminal.
// It blocks until the QR code is scanned or the context is cancelled.
func (s *Service) LoginWithQRCode(ctx context.Context, dbPath string) error {
	if err := s.initClient(ctx, dbPath); err != nil {
		return err
	}

	if s.client.Store.ID != nil {
		return s.client.Connect()
	}

	qrChan, _ := s.client.GetQRChannel(ctx)
	if err := s.client.Connect(); err != nil {
		return err
	}

	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		}
	}

	return nil
}

// LoginWithPairingCode authenticates using an 8-digit pairing code sent to phoneNumber.
// It blocks until the pairing code is generated; completion is signalled via WhatsApp events.
func (s *Service) LoginWithPairingCode(ctx context.Context, phoneNumber, dbPath string) (string, error) {
	if err := s.initClient(ctx, dbPath); err != nil {
		return "", err
	}

	if s.client.Store.ID != nil {
		return "", s.client.Connect()
	}

	if err := s.client.Connect(); err != nil {
		return "", err
	}

	code, err := s.client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		return "", fmt.Errorf("whatsapp: failed to get pairing code: %w", err)
	}

	return code, nil
}

// IsConnected reports whether the client is currently connected to WhatsApp.
func (s *Service) IsConnected() bool {
	return s.client != nil && s.client.IsConnected()
}

// Disconnect closes the connection to the WhatsApp servers.
func (s *Service) Disconnect() {
	if s.client != nil {
		s.client.Disconnect()
	}
}

// Send delivers subject and message to all configured recipients.
// The subject is rendered as bold text followed by a blank line before the message body.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	if !s.IsConnected() {
		return errors.New("whatsapp: client not connected, call LoginWithQRCode or LoginWithPairingCode first")
	}

	if len(s.recipients) == 0 {
		return ErrMissingRecipient
	}

	body := message
	if subject != "" {
		body = fmt.Sprintf("*%s*\n\n%s", subject, message)
	}

	for _, recipient := range s.recipients {
		_, err := s.client.SendMessage(ctx, recipient, &waE2E.Message{
			Conversation: proto.String(body),
		})
		if err != nil {
			return fmt.Errorf("whatsapp: failed to send to %s: %w", recipient, err)
		}
	}

	return nil
}
