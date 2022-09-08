package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type (
	// Service encapsulates a http daemon client.
	Service struct {
		client   *http.Client
		webhooks []string
	}

	// payload is the payload that is sent to the webhook. The content of this struct depends on the content-type of the
	// webhook.
	payload struct {
		Subject string `json:"subject,omitempty"`
		Message string `json:"message"`
	}
)

const (
	// TODO: Make these two configurable. Fixed for now to allows focus on the core logic.
	defaultContentType   = "application/json"
	defaultRequestMethod = "POST"
)

// New returns a new instance of a Service notification service. Parameter 'tag' is used as a log prefix and may be left
// empty, it has a fallback value.
func New() *Service {
	return &Service{
		client:   http.DefaultClient,
		webhooks: []string{},
	}
}

// AddReceivers takes one or more webhook URLs and appends them to the internal webhooks slice. The Send method will
// send a given message to all those webhooks. The webhook URLs must be valid and is expected to be an HTTP endpoint.
func (s *Service) AddReceivers(webhooks ...string) {
	s.webhooks = append(s.webhooks, webhooks...)
}

// WithClient sets the http client to be used for sending requests. Calling this method is optional, the default client
// will be used if this method is not called.
func (s *Service) WithClient(client *http.Client) {
	if client != nil {
		s.client = client
	}
}

// marshalPayload marshals the given payload into a reader. Currently, this is fixed to JSON. In the future, this method
// should be customizable based on the content-type of the webhook or be replaced by custom marshalling methods.
func (s *Service) marshalPayload(payload payload) (io.Reader, error) {
	// TODO: We need to allow different content types here. The content-type may vary depending on the webhook. Also,
	//       this probably needs to be controlled by the user, too.
	payloadRaw, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.Wrap(err, "marshal message")
	}

	return bytes.NewReader(payloadRaw), nil
}

// newRequest creates a new http request with the given method, content type, URL and payload.
func newRequest(ctx context.Context, method, contentType, url string, payload io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	return req, nil
}

// do sends the given request and returns an error if the request failed. A failed request gets identified by either
// an unsuccessful status code or a non-nil error.
func (s *Service) do(req *http.Request) error {
	// TODO: Pre-run middlewares

	// Actually send the HTTP request.
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// TODO: Post-run middlewares

	// Check if response code is 2xx. Should this be configurable?
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("responded with status code: %d", resp.StatusCode)
	}

	return nil
}

// send is a helper method that sends a message to a single webhook. It wraps the core logic of the Send method, which
// is creating a new request for the given webhook and sending it.
func (s *Service) send(ctx context.Context, webhook string, payload io.Reader) error {
	// Create a new HTTP request for the given webhook.
	req, err := newRequest(ctx, defaultRequestMethod, defaultContentType, webhook, payload)
	if err != nil {
		return errors.Wrapf(err, "create request for %q", webhook)
	}
	defer func() { _ = req.Body.Close() }()

	// Send the request
	err = s.do(req)
	if err != nil {
		return errors.Wrapf(err, "send request to %q", webhook)
	}

	return nil
}

// Send takes a message and sends it to all webhooks.
func (s *Service) Send(ctx context.Context, subject, message string) error {
	// Send message to all webhooks.
	for _, webhook := range s.webhooks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// TODO: At this point, we already need to know the content-type of the webhook. Marshalling has to be done
			//       based on the content-type. We need to allow the user to specify the content-type of the webhook and
			//       customize the marshalling process. A custom webhook type may be a good idea. This type could be
			//       responsible for marshalling and unmarshalling the payload and holding info about e.g. the request
			//       method and content-type.
			// Marshal the message into a payload.
			msg := payload{Subject: subject, Message: message}
			payload, err := s.marshalPayload(msg)
			if err != nil {
				return errors.Wrap(err, "marshal payload")
			}

			// Send the payload to the webhook.
			err = s.send(ctx, webhook, payload)
			if err != nil {
				return errors.Wrapf(err, "send message to %q", webhook)
			}
		}
	}

	return nil
}
