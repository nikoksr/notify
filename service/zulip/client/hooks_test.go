package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithBaseURL(t *testing.T) {
	t.Parallel()

	urls := []string{
		"https://localhost:5000",
		"https://zulip.com",
		"https://test.domain.co",
	}

	for i, url := range urls {
		client, _ := NewClient(WithBaseURL(url), WithCreds("<email>", "<apiKey>"))
		assert.Equal(
			t,
			client.baseURL,
			url,
			fmt.Sprintf("TEST %d: WithBaseURL hook failed", i),
		)
	}
}

func TestWithCreds(t *testing.T) {
	t.Parallel()

	creds := []struct {
		email  string
		apiKey string
	}{
		{email: "test@gmail.com", apiKey: "<apikey#1>"},
		{email: "zulip@gmail.com", apiKey: "<apikey#2>"},
		{email: "notify@gmail.com", apiKey: "<apikey#3>"},
	}

	for i, cred := range creds {
		client, err := NewClient(WithCreds(cred.email, cred.apiKey))
		if err != nil {
			t.Errorf("TEST %d: WithCreds hook errored", i)
		} else {
			assert.Equal(
				t,
				client.email,
				cred.email,
				fmt.Sprintf("TEST %d: WithCreds hook failed for email", i),
			)
			assert.Equal(
				t,
				client.apiKey,
				cred.apiKey,
				fmt.Sprintf("TEST %d: WithCreds hook failed for apiKey", i),
			)
		}
	}
}

func TestWithTimeout(t *testing.T) {
	t.Parallel()

	timeouts := []time.Duration{
		1 * time.Minute,
		50 * time.Second,
		200 * time.Millisecond,
	}

	for i, timeout := range timeouts {
		client, _ := NewClient(WithTimeout(timeout), WithCreds("<email>", "<apiKey>"))

		assert.Equal(
			t,
			client.timeout,
			timeout,
			fmt.Sprintf("TEST %d: WithTimeout hook failed", i),
		)
	}
}
