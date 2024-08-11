package homeassistant

import (
	"context"
	stdhttp "net/http"
	"strings"

	"github.com/nikoksr/notify/service/http"
)

// HomeAssistant struct holds necessary data to communicate with the HomeAssistant API.
type HomeAssistant struct {
	httpClient http.Service
}

// New returns a new instance of a HomeAssistant notification service.
func New() (*HomeAssistant, error) {
	h := HomeAssistant{
		httpClient: *http.New(),
	}
	return &h, nil
}

// AddWebhook takes HomeAssistant automation webhooks and adds them. The Send method will send to all these
// baseUrl should be the home assistant base url for example https://your-home-assistant:8123
// hookId should be what was set at the trigger of the webhook
// method should be the method selected at the trigger of the webhook
// For more information read HomeAssistant documentation at:
// - https://www.home-assistant.io/docs/automation/trigger/#webhook-trigger
func (h *HomeAssistant) AddWebhook(baseUrl string, hookId string, method string) {
	baseUrl = strings.TrimRight(baseUrl, "/")
	u := baseUrl + "/api/webhook/" + hookId

	// From example:
	// curl -X POST -d 'key=value&key2=value2' https://your-home-assistant:8123/api/webhook/some_hook_id

	hook := &http.Webhook{
		URL:         u,
		Header:      stdhttp.Header{},
		ContentType: "application/json",
		Method:      strings.ToUpper(method),
		BuildPayload: func(subject, message string) (payload any) {
			dataMap := make(map[string]string)
			dataMap["subject"] = subject
			dataMap["message"] = message

			return dataMap
		},
	}

	h.httpClient.AddReceivers(hook)
}

// Send takes a subject and a message and sends them to all previously set webhooks
func (h HomeAssistant) Send(ctx context.Context, subject, message string) error {
	// TODO: should setup mqtt automation integration as well
	//       https://www.home-assistant.io/docs/automation/trigger/#mqtt-trigger

	return h.httpClient.Send(ctx, subject, message)
}
