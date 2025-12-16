package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

//go:generate mockery --name=GotifyService --output=. --case=underscore --inpackage
type GotifyService interface {
	Send(ctx context.Context, subject, message string) error
}

const DefaultPriority = 1

var _ GotifyService = &Gotify{}

// Gotify struct holds necessary data to communicate with Gotify API
type Gotify struct {
	httpClient *http.Client
	baseURL    string
	appToken   string
	priority   int
}

// New returns a new instance of a Gotify notification service.
// For more information about Gotify api credential:
//
//	-> https://gotify.net/docs/pushmsg
func New(token string, baseUrl string) *Gotify {
	return &Gotify{
		httpClient: http.DefaultClient,
		baseURL:    baseUrl,
		appToken:   token,
	}
}

// NewWithPriority returns a new instance of a Gotify notification service with custom priority.
// For more information about Gotify api credential:
//
//	-> https://gotify.net/docs/pushmsg
func NewWithPriority(token string, baseUrl string, priority int) *Gotify {
	return &Gotify{
		httpClient: http.DefaultClient,
		baseURL:    baseUrl,
		appToken:   token,
		priority:   priority,
	}
}

type Message struct {
	Title    string
	Message  string
	Priority int
}

// Send the message subject and message body to the Gotify service and set the priority.
func (gotify *Gotify) Send(ctx context.Context, subject, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		reqBody := &Message{
			Title:    subject,
			Message:  message,
			Priority: DefaultPriority,
		}

		// Set the default priority when the user does not set a priority.
		if gotify.priority != 0 {
			reqBody.Priority = gotify.priority
		}

		// Make request
		body, err := json.Marshal(reqBody)
		if err != nil {
			return errors.Wrap(err, "encode message body")
		}

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/message", gotify.baseURL),
			bytes.NewReader(body),
		)
		if err != nil {
			return errors.Wrap(err, "create new request")
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", gotify.appToken))
		req.Header.Set("Content-Type", "application/json; charset=utf-8")

		// Send request to gotify service
		resp, err := gotify.httpClient.Do(req)
		if err != nil {
			return errors.Wrapf(err, "send request to gotify server")
		}

		// Read response and verify success
		result, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read response")
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("gotify returned status code %d: %s", resp.StatusCode, string(result))
		}

		return nil
	}
}
