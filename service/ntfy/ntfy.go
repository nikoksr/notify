package ntfy

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Ntfy struct holds necessary address config to send push notification.
type Ntfy struct {
	ntfyAddress string
	ntfyMethod  string
}

type BodyOptions struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message"`
	Title    string   `json:"title"`
	Tags     []string `json:"tags"`
	Priority int      `json:"priority"`
	Attach   string   `json:"attach"`
	Filename string   `json:"filename"`
	Click    string   `json:"click"`
	Actions  []struct {
		Action string `json:"action"`
		Label  string `json:"label"`
		URL    string `json:"url"`
	} `json:"actions"`
}

// New returns a new instance of a Ntfy notification service.
func New() *Ntfy {
	return &Ntfy{
		ntfyAddress: "https://ntfy.sh",
		ntfyMethod:  "POST",
	}
}

func (m *Ntfy) Send(ctx context.Context, topic string, body string) error {

	bodyByte := []byte(body)

	isJson := json.Valid(bodyByte)
	if !isJson {
		err := errors.New("Invalid JSON Body")
		return err
	}

	var bodyj BodyOptions
	if err := json.Unmarshal(bodyByte, &bodyj); err != nil {
		err := errors.New("topic is not defined")
		return err
	}

	if bodyj.Topic == "" {
		bodyj.Topic = topic
	}

	newBody, _ := json.Marshal(bodyj)

	req, _ := http.NewRequest(m.ntfyMethod, m.ntfyAddress, strings.NewReader(string(newBody)))
	_, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrap(err, "failed to send push notification")
		return err
	}

	return nil
}
