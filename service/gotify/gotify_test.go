package gotify

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGotify_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	assert.NotNil(New(
		"testingAppToken",
		"baseUrl",
	))
}

func TestGotify_SendSuccess(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	ctx := context.Background()

	mock := newMockGotifyService(t)
	mock.On("Send", ctx, "testing", "hello world!").Return(nil)

	err := mock.Send(ctx, "testing", "hello world!")
	assert.Nil(err)

	mock.AssertCalled(t, "Send", ctx, "testing", "hello world!")
}

func TestGotify_SendFailure(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	ctx := context.Background()

	mock := newMockGotifyService(t)
	mock.On("Send", ctx, "testing", "hello world!").Return(errors.New("failed to send message"))

	err := mock.Send(ctx, "testing", "hello world!")
	assert.NotNil(err)

	mock.AssertCalled(t, "Send", ctx, "testing", "hello world!")
}
