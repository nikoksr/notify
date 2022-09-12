package fcm

import (
	"context"

	"github.com/appleboy/go-fcm"
	"github.com/pkg/errors"
)

// Compile-time check that fcm.Client satisfies fcmClient interface.
var _ fcmClient = &fcm.Client{}

var (
	// DataKey is used as a context.Context key to optionally add data to the message payload.
	DataKey = msgDataKey{}
	// RetriesKey is used as a context.Context key to optionally set a total of retry attempts per each message.
	RetriesKey = msgRetriesKey{}
)

type (
	msgDataKey    struct{}
	msgRetriesKey struct{}
)

// fcmClient abstracts go-fcm for writing unit tests
//
//go:generate mockery --name=fcmClient --output=. --case=underscore --inpackage
type fcmClient interface {
	SendWithRetry(*fcm.Message, int) (*fcm.Response, error)
}

// Service encapsulates the FCM client along with internal state for storing device tokens.
type Service struct {
	client       fcmClient
	deviceTokens []string
}

// New returns a new instance of a FCM notification service.
func New(serverAPIKey string) (*Service, error) {
	client, err := fcm.NewClient(serverAPIKey)
	if err != nil {
		return nil, err
	}

	s := &Service{
		client:       client,
		deviceTokens: []string{},
	}
	return s, nil
}

// AddReceivers takes FCM device tokens and appends them to the internal device tokens slice.
// The Send method will send a given message to all those devices.
func (s *Service) AddReceivers(deviceTokens ...string) {
	s.deviceTokens = append(s.deviceTokens, deviceTokens...)
}

// Send takes a message subject and a message body and sends them to all previously set devices.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	msg := &fcm.Message{
		Notification: &fcm.Notification{
			Title: subject,
			Body:  message,
		},
	}

	if data, ok := getMessageData(ctx); ok {
		msg.Data = data
	}

	retryAttempts := getMessageRetryAttempts(ctx)

	for _, deviceToken := range s.deviceTokens {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg.To = deviceToken

			_, err := s.client.SendWithRetry(msg, retryAttempts)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to FCM device with token '%s'", deviceToken)
			}
		}
	}

	return nil
}

func getMessageData(ctx context.Context) (data map[string]interface{}, ok bool) {
	value := ctx.Value(DataKey)
	if value != nil {
		data, ok = value.(map[string]interface{})
	}
	return
}

func getMessageRetryAttempts(ctx context.Context) int {
	value := ctx.Value(RetriesKey)
	if value != nil {
		if retryAttempts, ok := value.(int); ok {
			return retryAttempts
		}
	}
	return 0
}
