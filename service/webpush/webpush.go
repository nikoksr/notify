package webpush

import (
	"context"
	"encoding/json"
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/pkg/errors"
)

var (
	// DataKey is used as a context.Context key to optionally add data to the message payload.
	DataKey = msgDataKey{}
	// UrgencyKey is used as context.Context key to optionally add urgency to the message.
	UrgencyKey = msgUrgencyKey{}
)

type (
	msgDataKey    struct{}
	msgUrgencyKey struct{}
)

type Service struct {
	subscriptions   []webpush.Subscription
	VAPIDPublicKey  string
	VAPIDPrivateKey string
}

func New(VAPIDPublicKey string, VAPIDPrivateKey string) (*Service, error) {
	subscriptions := []webpush.Subscription{}

	return &Service{
		subscriptions,
		VAPIDPublicKey,
		VAPIDPrivateKey,
	}, nil
}

func (s *Service) AddReceivers(subscriptions ...string) error {
	// Parses subscription objects and saves them into the service
	for _, subJSON := range subscriptions {
		subscription := webpush.Subscription{}
		err := json.Unmarshal([]byte(subJSON), &subscription)
		if err != nil {
			return err
		}
		s.subscriptions = append(s.subscriptions, subscription)
	}

	return nil
}

func (s *Service) Send(ctx context.Context, subject, message string) error {
	options := getOptionsFromCtx(ctx)
	options.VAPIDPrivateKey = s.VAPIDPrivateKey
	options.VAPIDPublicKey = s.VAPIDPublicKey

	payload, err := getMessageFromCtx(ctx, subject, message)
	if err != nil {
		return err
	}

	for i := range s.subscriptions {
		sub := s.subscriptions[i]
		res, err := webpush.SendNotificationWithContext(ctx, payload, &sub, &options)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to webpush subscription %+v", sub)
		}
		if res.StatusCode != 201 {
			return errors.Wrapf(
				fmt.Errorf("expected StatusCode 201, got %v", res.StatusCode),
				"failed to send message to webpush subscription %+v", sub,
			)
		}
	}

	return nil
}

func getOptionsFromCtx(ctx context.Context) webpush.Options {
	return webpush.Options{
		Urgency: getUrgencyFromCtx(ctx),
	}
}

// Gets Urgency From Ctx, returns webpush.UrgencyNormal as a default
func getUrgencyFromCtx(ctx context.Context) webpush.Urgency {
	value := ctx.Value(UrgencyKey)
	if value == nil {
		return webpush.UrgencyNormal
	}
	data, ok := value.(webpush.Urgency)

	if !ok {
		return webpush.UrgencyNormal
	}

	// validate urgency
	switch data {
	case webpush.UrgencyVeryLow, webpush.UrgencyLow, webpush.UrgencyNormal, webpush.UrgencyHigh:
		return data
	}

	return webpush.UrgencyNormal
}

type webpushMessage struct {
	Subject string                 `json:"subject"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func getMessageFromCtx(ctx context.Context, subject, message string) ([]byte, error) {
	webpushMsg := webpushMessage{
		Subject: subject,
		Message: message,
		Data:    map[string]interface{}{},
	}
	data, ok := getMessageData(ctx)
	if ok {
		webpushMsg.Data = data
	}
	payload, err := json.Marshal(webpushMsg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize message")
	}

	return payload, nil
}

func getMessageData(ctx context.Context) (data map[string]interface{}, ok bool) {
	value := ctx.Value(DataKey)
	if value != nil {
		data, ok = value.(map[string]interface{})
	}
	return
}
