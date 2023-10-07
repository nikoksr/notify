package apns2

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	apnsSvc "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

// Compile-time check that fcm.Client satisfies fcmClient interface.
var _ apns2Client = &apnsSvc.Client{}

//go:generate mockery --name=apns2Client --output=. --case=underscore --inpackage
type apns2Client interface {
	Push(n *apnsSvc.Notification) (*apnsSvc.Response, error)
}

// hook to parse p12 bytes for credentials
func P12Bytes(bytes []byte, password string) func() (apns2Client, error) {
	return func() (apns2Client, error) {
		cert, err := certificate.FromP12Bytes(bytes, password)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid certificates %s %s", bytes, password)
		}

		client := apnsSvc.NewClient(cert).Production()
		return client, nil
	}
}

// hook to parse p12 file for credentials
func P12File(filename, password string) func() (apns2Client, error) {
	return func() (apns2Client, error) {
		cert, err := certificate.FromP12File(filename, password)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid certificates %s %s", filename, password)
		}

		client := apnsSvc.NewClient(cert).Production()
		return client, nil
	}
}

// hook to parse pem file for credentials
func PemFile(filename, password string) func() (apns2Client, error) {
	return func() (apns2Client, error) {
		cert, err := certificate.FromPemFile(filename, password)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid certificates %s %s", filename, password)
		}

		client := apnsSvc.NewClient(cert).Production()
		return client, nil
	}
}

// hook to parse pem bytes for credentials
func PemBytes(bytes []byte, password string) func() (apns2Client, error) {
	return func() (apns2Client, error) {
		cert, err := certificate.FromPemBytes(bytes, password)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid certificates %s %s", bytes, password)
		}

		client := apnsSvc.NewClient(cert).Production()
		return client, nil
	}
}

func buildNotification(token, topic, msg string) *apnsSvc.Notification {
	notification := &apnsSvc.Notification{}
	notification.DeviceToken = token
	notification.Topic = topic
	notification.Payload = []byte(fmt.Sprintf(`{"aps":{"alert":"%s"}}`, msg))

	return notification
}

// Service encapsulates the APNS2 client along with internal state for storing device tokens.
type Service struct {
	client       apns2Client
	topic        string
	deviceTokens []string
}

// New returns a new instance of a APNS2 notification service
func New(makeClient func() (apns2Client, error), topic string) (*Service, error) {
	apnsClient, err := makeClient()
	if err != nil {
		return nil, err
	}

	client := &Service{
		apnsClient,
		topic,
		make([]string, 0),
	}
	return client, nil
}

// AddReceivers takes APNS2 device tokens and appends them to the internal device tokens slice.
// The Send method will send a given message to all those devices.
func (s *Service) AddReceivers(deviceTokens ...string) {
	s.deviceTokens = append(s.deviceTokens, deviceTokens...)
}

// Send takes a message subject and a message body and sends them to all previously set devices.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	for _, deviceToken := range s.deviceTokens {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			notification := buildNotification(deviceToken, s.topic, subject+" "+message)

			_, err := s.client.Push(notification)
			if err != nil {
				return errors.Wrapf(err, "failed to send notification to %s", deviceToken)
			}
		}
	}

	return nil
}
