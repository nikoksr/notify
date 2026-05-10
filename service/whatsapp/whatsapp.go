package whatsapp

import (
    "context"
    "errors"
    "fmt"
    "os"

    "github.com/mdp/qrterminal/v3"
    "go.mau.fi/whatsmeow"
    "go.mau.fi/whatsmeow/proto/wa"
    "go.mau.fi/whatsmeow/store/sqlstore"
    "go.mau.fi/whatsmeow/types"
    waLog "go.mau.fi/whatsmeow/util/log"
    "google.golang.org/protobuf/proto"
    _ "modernc.org/sqlite" // Import pure Go SQLite driver for whatsmeow store.
)

var ErrMissingRecipient = errors.New("whatsapp: missing recipient JID")

// Service implements notify.Notifier for WhatsApp using whatsmeow.
type Service struct {
    client     *whatsmeow.Client
    recipients []types.JID
    logger     waLog.Logger
    dbPath     string
}

// New returns a new instance of a WhatsApp notification service.
func New() (*Service, error) {
    return &Service{
        logger: waLog.Noop,
    }, nil
}

// AddReceivers takes WhatsApp contacts and adds them to the internal contacts list.
// Format: 6281234567890@s.whatsapp.net.
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
    s.dbPath = dbPath
    dbLog := waLog.Stdout("Database", "INFO", true)
    container, err := sqlstore.New(ctx, "sqlite", "file:"+dbPath+"?_foreign_keys=on", dbLog)
    if err != nil {
        return fmt.Errorf("whatsapp: failed to init store: %w", err)
    }

    deviceStore, err := container.GetFirstDevice(ctx)
    if err != nil {
        return fmt.Errorf("whatsapp: failed to get device: %w", err)
    }

    s.client = whatsmeow.NewClient(deviceStore, s.logger)
    return nil
}

// LoginWithQRCode authenticates using QR code on terminal.
// It will block until QR is scanned or context timeout.
func (s *Service) LoginWithQRCode(ctx context.Context, dbPath string) error {
    if err := s.initClient(ctx, dbPath); err != nil {
        return err
    }

    if s.client.Store.ID != nil {
        if err := s.client.Connect(); err != nil {
            return err
        }
        s.logger.Infof("Already logged in, auto connected")
        return nil
    }

    qrChan, _ := s.client.GetQRChannel(ctx)
    if err := s.client.Connect(); err != nil {
        return err
    }

    for evt := range qrChan {
        if evt.Event == "code" {
            qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
        } else {
            s.logger.Infof("Login event: %s", evt.Event)
        }
    }
    return nil
}

// LoginWithPairingCode authenticates using an 8-digit pairing code.
// It will block until the code is generated; pairing completion is handled via events.
func (s *Service) LoginWithPairingCode(ctx context.Context, phoneNumber, dbPath string) error {
    if err := s.initClient(ctx, dbPath); err != nil {
        return err
    }

    if s.client.Store.ID != nil {
        s.logger.Infof("Already logged in, auto connected")
        return s.client.Connect()
    }

    if err := s.client.Connect(); err != nil {
        return err
    }

    code, err := s.client.PairPhone(ctx, phoneNumber, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
    if err != nil {
        return fmt.Errorf("whatsapp: failed to get pairing code: %w", err)
    }

    s.logger.Infof("Your Pairing Code: %s", code)
    return nil
}

// Disconnect closes the connection to the WhatsApp servers.
func (s *Service) Disconnect() {
    if s.client != nil {
        s.client.Disconnect()
    }
}

// Send takes a message subject and a message body and sends them to all previously set contacts.
// Subject will be formatted as bold text.
func (s *Service) Send(ctx context.Context, subject, message string) error {
    if s.client == nil || !s.client.IsConnected() {
        return errors.New("whatsapp: client not connected, call LoginWithQRCode() or LoginWithPairingCode() first")
    }
    if len(s.recipients) == 0 {
        return ErrMissingRecipient
    }

    fullMsg := message
    if subject != "" {
        fullMsg = fmt.Sprintf("*%s*\n\n%s", subject, message)
    }

    for _, recipient := range s.recipients {
        _, err := s.client.SendMessage(ctx, recipient, &wa.Message{
            Conversation: proto.String(fullMsg),
        })
        if err != nil {
            return fmt.Errorf("whatsapp: failed to send to %s: %w", recipient, err)
        }
    }
    return nil
}
