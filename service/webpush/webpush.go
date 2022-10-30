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

// These are exposed Urgency constants from the webpush package.
var (
	// UrgencyVeryLow requires device state: on power and Wi-Fi
  UrgencyVeryLow = webpush.UrgencyVeryLow
	// UrgencyLow requires device state: on either power or Wi-Fi
  UrgencyLow = webpush.UrgencyLow
	// UrgencyNormal excludes device state: low battery
  UrgencyNormal = webpush.UrgencyNormal
	// UrgencyHigh admits device state: low battery
  UrgencyHigh = webpush.UrgencyHigh
)

type (
	msgDataKey    struct{}
	msgUrgencyKey struct{}
)

// Service encapsulates the webpush notification system along with the internal state
type Service struct {
	subscriptions   []webpush.Subscription
	vapidPublicKey  string
	vapidPrivateKey string
}

// New returns a new instance of the Service
func New(vapidPublicKey string, vapidPrivateKey string) *Service {
	subscriptions := []webpush.Subscription{}

	return &Service{
		subscriptions,
		vapidPublicKey,
		vapidPrivateKey,
	}
}

// AddReceivers accepts multiple JSON strings as []byte, each representing a subscription (refer to: https://developer.mozilla.org/en-US/docs/Web/API/PushSubscription)
// Send method will send the notification to all these subscriptions
func (s *Service) AddReceivers(subscriptions ...[]byte) error {
	// Parses subscription objects and saves them into the service
	for _, subJSON := range subscriptions {
		subscription := webpush.Subscription{}
		err := json.Unmarshal(subJSON, &subscription)
		if err != nil {
			return err
		}
		s.subscriptions = append(s.subscriptions, subscription)
	}

	return nil
}

// Send sends the message to all the subscriptions that were previously added
// It accepts the following options from the ctx 
//  * DataKey - This is a map[string]interface{} which sent as extra payload (default: empty)
//  * Urgency - This is a webpush.Urgency or string (default: UrgencyNormal)
// All these are optinal and are have sensible defaults in place
func (s *Service) Send(ctx context.Context, subject, message string) error {
	options := getOptionsFromCtx(ctx)
	options.VAPIDPrivateKey = s.vapidPrivateKey
	options.VAPIDPublicKey = s.vapidPublicKey

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

// getOptionsFromCtx extracts webpush.Options from the given ctx
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


// webpushMessage  a webpush message that is serialized into JSON
type webpushMessage struct {
	Subject string                 `json:"subject"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

// getMessageFromCtx constructs a webpushMessage from the ctx, subject and message.
// It returns the webpushMessage converted to JSON []byte 
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

// getMessageData extracts message data out of the context
func getMessageData(ctx context.Context) (data map[string]interface{}, ok bool) {
	value := ctx.Value(DataKey)
	if value != nil {
		data, ok = value.(map[string]interface{})
	}
	return
}
