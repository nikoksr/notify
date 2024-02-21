package zulip

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	gzb "github.com/ifo/gozulipbot"
	"github.com/pkg/errors"
)

//go:generate mockery --name=zulipClient --output=. --case=underscore --inpackage
type zulipClient interface {
	Message(gzb.Message) (*http.Response, error)
}

// Compile-time check to ensure that zulip message client implements the zulipClient interface.
var _ zulipClient = new(gzb.Bot)

// Zulip struct holds necessary data to communicate with the Zulip API.
type Zulip struct {
	client    zulipClient
	receivers []*Receiver
}

func New(domain, apiKey, botEmail string) *Zulip {
	client := &gzb.Bot{
		APIURL: fmt.Sprintf("https://%s.zulipchat.com/api/v1/", domain),
		APIKey: apiKey,
		Email:  botEmail,
	}

	client.Init()

	zulip := &Zulip{
		client:    client,
		receivers: make([]*Receiver, 0),
	}

	return zulip
}

func (z *Zulip) AddReceivers(receivers ...*Receiver) {
	z.receivers = append(z.receivers, receivers...)
}

func (z *Zulip) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	for _, receiver := range z.receivers {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			emails := make([]string, 0)
			if receiver.email != "" {
				emails = append(emails, receiver.email)
			}

			msg := gzb.Message{
				Content: fullMessage,
				Emails:  emails,
				Stream:  receiver.stream,
				Topic:   receiver.topic,
			}

			resp, err := z.client.Message(msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Zulip receiver")
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)

			switch resp.StatusCode {
			case http.StatusBadRequest:
				var errorResp ErrorResponse
				_ = json.Unmarshal(body, &errorResp)

				return errors.Errorf("failed to send message to Zulip receiver: %s", errorResp.Message)

			case http.StatusOK:
				break

			default:
				return errors.Errorf("failed to send message to Zulip receiver: %s", body)
			}
		}
	}

	return nil
}
