package bark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Service struct {
	Server    string
	DeviceKey string
	url       string
}

const DefaultServer = "api.day.app"

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

		resp, err := http.Post(p.url, "application/json; charset=utf-8", bytes.NewReader(data))
		if err != nil {
			return errors.Wrap(err, "send bark request failed")
		}
		result, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

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
