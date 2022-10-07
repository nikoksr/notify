package metaworkplace

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// ENDPOINT is the base URL of the Workplace API to send messages.
	ENDPOINT = "https://graph.facebook.com/me/messages"
)

// metaWorkplaceService is the internal implementation of the Meta Workplace notification service.
//
//go:generate mockery --name=metaWorkplaceService --output=. --case=underscore --inpackage
type metaWorkplaceService interface {
	send(payload interface{}) *MetaWorkplaceResponse
}

// ValidateConfig checks if the required configuration data is present.
func (sc *MetaWorkplaceServiceConfig) ValidateConfig() error {
	if sc.AccessToken == "" {
		return errors.New("a valid Meta Workplace access token is required")
	}
	return nil
}

// New returns a new instance of a Meta Workplace notification service.
func New(token string) *MetaWorkplaceService {
	serviceConfig := MetaWorkplaceServiceConfig{
		AccessToken: token,
		Endpoint:    ENDPOINT,
	}

	err := serviceConfig.ValidateConfig()
	if err != nil {
		log.Fatalf("failed to validate Meta Workplace service configuration: %v", err)
	}

	return &MetaWorkplaceService{
		MetaWorkplaceServiceConfig: serviceConfig,
		userIDs:                    []string{},
		threadIDs:                  []string{},
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// AddThreads takes Workplace thread IDs and adds them to the internal thread ID list. The Send method will send
// a given message to all those threads.
func (mw *MetaWorkplaceService) AddThreads(threadIDs ...string) {
	mw.threadIDs = append(mw.threadIDs, threadIDs...)
}

// AddUsers takes user IDs and adds them to the internal user ID list. The Send method will send
// a given message to all those users.
func (mw *MetaWorkplaceService) AddUsers(userIDs ...string) {
	mw.userIDs = append(mw.userIDs, userIDs...)
}

// send takes a payload, sends it to the Workplace API, and returns the response.
func (mw *MetaWorkplaceService) send(payload interface{}) *MetaWorkplaceResponse {
	data, err := json.Marshal(payload)
	if err != nil {
		response := MetaWorkplaceResponse{
			MessageID: "",
			ThreadKey: "",
			Error: &MetaWorkplaceErrorResponse{
				Message: "failed to marshal payload",
			},
		}
		return &response
	}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest(http.MethodPost, mw.Endpoint, buff)
	if err != nil {
		log.Println("failed to create new HTTP request")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+mw.AccessToken)

	res, err := mw.client.Do(req)
	if err != nil {
		log.Printf("failed to send HTTP request: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(res.Body)

	data, err = io.ReadAll(res.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
	}

	var response MetaWorkplaceResponse

	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Printf("failed to unmarshal response body: %v", err)
	}

	return &response
}

// Send takes a message and sends it to the provided user and/or thread IDs.
func (mw *MetaWorkplaceService) Send(ctx context.Context, subject string, message string) error {
	if len(mw.userIDs) == 0 && len(mw.threadIDs) == 0 {
		return errors.New("no user or thread IDs provided")
	}

	if len(mw.threadIDs) != 0 {
		for _, threadID := range mw.threadIDs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				payload := metaWorkplaceThreadPayload{
					Message:   metaWorkplaceMessage{Text: message},
					Recipient: metaWorkplaceThread{ThreadID: threadID},
				}
				err := mw.send(payload)
				if err.Error != nil {
					log.Printf("%+v\n", err.Error)
					return errors.New("failed to send message to Workplace thread: " + threadID)
				}
			}
		}
	}

	if len(mw.userIDs) != 0 {
		for _, userID := range mw.userIDs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				payload := metaWorkplaceUserPayload{
					Message:   metaWorkplaceMessage{Text: message},
					Recipient: metaWorkplaceUsers{UserIDs: []string{userID}},
				}
				err := mw.send(payload)
				if err.Error != nil {
					log.Printf("%+v\n", err.Error)
					return errors.New("failed to send message to Workplace user: " + userID)
				}

				// Capture the thread ID for the user and add it to the thread ID list. Subsequent
				// messages will be sent to the thread instead of creating a new thread for the user.
				mw.threadIDs = append(mw.threadIDs, err.ThreadKey)
			}
		}

		// Clear the user ID list. Subsequent messages will be sent to the thread instead of creating a new thread for the user.
		mw.userIDs = []string{}
	}

	return nil
}
