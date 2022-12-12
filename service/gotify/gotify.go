package gotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

//go:generate mockery --name=gotifyService --output=. --case=underscore --inpackage
type gotifyService interface {
	Send(ctx context.Context, subject, message string) error
}

var _ gotifyService = &Gotify{}

// Gotify struct holds necessary data to communicate with Gotify API
type Gotify struct {
	httpClient *http.Client
	baseUrl    string
	appToken   string
}

func New(appToken, baseUrl string) *Gotify {
	g := &Gotify{
		httpClient: http.DefaultClient,
		baseUrl:    baseUrl,
		appToken:   appToken,
	}

	return g
}

type newMessageRequestBody struct {
	Title    string
	Message  string
	Priority int
}

func (g *Gotify) Send(ctx context.Context, subject, message string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		reqBody := &newMessageRequestBody{
			Title:    subject,
			Message:  message,
			Priority: 1,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/message", g.baseUrl), bytes.NewReader(body))
		if err != nil {
			return err
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", g.appToken))

		_, err = g.httpClient.Do(req)
		if err != nil {
			return errors.Wrapf(err, "failed to send message to gotify server")
		}

		return nil
	}
}
