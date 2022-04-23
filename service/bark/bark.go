package bark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// Service allow you to configure Bark service.
type Service struct {
	Server    string
	DeviceKey string
	url       string
}

// DefaultServer is the default server to use for the bark service.
const DefaultServer = "api.day.app"

// New returns a new instance of Bark service.
func New(deviceKey, server string) *Service {
	p := &Service{
		Server:    server,
		DeviceKey: deviceKey,
	}
	if server != "" {
		p.url = fmt.Sprintf("https://%s/push", server)
	} else {
		p.url = fmt.Sprintf("https://%s/push", "api.day.app")
	}
	return p
}

// Send takes a message subject and a message content and sends them to bark application.
func (p *Service) Send(ctx context.Context, subject, content string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		type postData struct {
			DeviceKey string `json:"device_key"`
			Title     string `json:"title"`
			Body      string `json:"body,omitempty"`
			Badge     int    `json:"badge,omitempty"`
			Sound     string `json:"sound,omitempty"`
			Icon      string `json:"icon,omitempty"`
			Group     string `json:"group,omitempty"`
			URL       string `json:"url,omitempty"`
		}
		pd := &postData{
			DeviceKey: p.DeviceKey,
			Title:     subject,
			Body:      content,
			Sound:     "alarm.caf",
		}
		data, err := json.Marshal(pd)
		if err != nil {
			return err
		}
		req, err := http.NewRequestWithContext(ctx, "POST", p.url, bytes.NewBuffer(data))
		if err != nil {
			return errors.Wrap(err, "failed to create bark request")
		}

		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		httpClient := &http.Client{
			Timeout: time.Second * 5,
		}
		resp, err := httpClient.Do(req)
		if err != nil {
			return errors.Wrap(err, "send bark request failed")
		}
		defer resp.Body.Close()

		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("statusCode: %d, body: %v", resp.StatusCode, string(result))
			return errors.Wrap(err, "send bark message failed")
		}
		if err != nil {
			return errors.Wrapf(err, "failed to send message")
		}
		return nil
	}
}
