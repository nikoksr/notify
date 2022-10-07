package metaworkplace

import (
	"net/http"
)

type (
	// MetaWorkplaceErrorResponse is a custom error type for the Meta Workplace service. This struct should be filled when
	// a status code other than 200 is received from the Meta Workplace API.
	MetaWorkplaceErrorResponse struct {
		Message      string `json:"message"`
		Type         string `json:"type"`
		Code         int    `json:"code"`
		ErrorSubcode int    `json:"error_subcode"`
		FbtraceID    string `json:"fbtrace_id"`
	}

	// MetaWorkplaceResponse is a custom response type for the Meta Workplace service. Only the MessageID and
	// ThreadKey fields should be filled when a 200 response is received from the Meta Workplace API.
	MetaWorkplaceResponse struct {
		MessageID string                      `json:"message_id"`
		ThreadKey string                      `json:"thread_key"`
		Error     *MetaWorkplaceErrorResponse `json:"error"`
	}

	// MetaWorkplaceService struct holds necessary data to communicate with Meta Workplace users and/or threads.
	MetaWorkplaceService struct {
		MetaWorkplaceServiceConfig
		userIDs   []string
		threadIDs []string
		client    *http.Client
	}

	// MetaWorkplaceServiceConfig holds the required configuration data when creating a connection to the Meta
	// Workplace API.
	MetaWorkplaceServiceConfig struct {
		AccessToken string
		Endpoint    string
	}

	// metaWorkplaceUserPayload is a custom payload type for the Meta Workplace service. This struct should be filled
	// when sending a message to a user.
	metaWorkplaceUserPayload struct {
		Message   metaWorkplaceMessage `json:"message"`
		Recipient metaWorkplaceUsers   `json:"recipient"`
	}

	// metaWorkplaceUserPayload is a custom payload type for the Meta Workplace service. This struct should be filled
	// when sending a message to a thread.
	metaWorkplaceThreadPayload struct {
		Message   metaWorkplaceMessage `json:"message"`
		Recipient metaWorkplaceThread  `json:"recipient"`
	}

	// metaWorkplaceUsers holds the user IDs to send a message to.
	metaWorkplaceUsers struct {
		UserIDs []string `json:"ids"`
	}

	// metaWorkplaceThread holds the thread ID to send a message to.
	metaWorkplaceThread struct {
		ThreadID string `json:"thread_key"`
	}

	// metaWorkplaceMessage holds the message to send.
	metaWorkplaceMessage struct {
		Text string `json:"text"`
	}
)
