package mastodon

import (
	"context"
	"errors"
	"testing"

	"github.com/mattn/go-mastodon"
	"github.com/stretchr/testify/assert"
)

var errSome = errors.New("some error")

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		cfg           *Config
		expectedError error
	}{
		{
			name: "Success",
			cfg: &Config{
				ServerURL:   "some-url",
				AccessToken: "some-access-token",
			},
			expectedError: nil,
		},
		{
			name: "Invalid configuration",
			cfg: &Config{
				ServerURL:   "",
				AccessToken: "",
			},
			expectedError: ErrInvalidConfParams,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.cfg)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestSend(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		setup         func(*mockmastodonClient)
		expectedError error
	}{
		{
			name:    "Success",
			message: "Hello world!",
			setup: func(mc *mockmastodonClient) {
				mc.On("PostStatus", context.Background(), &mastodon.Toot{
					Status: "Hello world!",
				}).
					Return(&mastodon.Status{}, nil)
			},
			expectedError: nil,
		},
		{
			name:    "Sending failed",
			message: "Hello world!",
			setup: func(mc *mockmastodonClient) {
				mc.On("PostStatus", context.Background(), &mastodon.Toot{
					Status: "Hello world!",
				}).
					Return(nil, errSome)
			},
			expectedError: nil,
		},
	}

	client := new(mockmastodonClient)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup(client)

			mastodon := Mastodon{
				client: client,
			}

			err := mastodon.Send(context.Background(), tc.message, tc.message)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
