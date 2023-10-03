package client

import (
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
		if assert.Equal(t, client.baseURL, url) {
			t.Logf("TEST %d: WithBaseURL hook passed", i)
		} else {
			t.Errorf("TEST %d: WithBaseURL hook failed", i)
		}
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
			if assert.Equal(t, client.email, cred.email) && assert.Equal(t, client.apiKey, cred.apiKey) {
				t.Logf("TEST %d: WithCreds hook passed", i)
			} else {
				t.Errorf("TEST %d: WithCreds hook failed", i)
			}
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

		if assert.Equal(t, client.timeout, timeout) {
			t.Logf("TEST %d: WithTimeout hook passed", i)
		} else {
			t.Errorf("TEST %d: WithTimeout hook failed", i)
		}
	}
}
