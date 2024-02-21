package zulip

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/require"
)

func TestZulip_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
	assert.NotNil(service)
}

func TestZulip_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
	assert.NotNil(service)

	service.AddReceivers(Direct("some-user@email.com"))
	assert.Len(service.receivers, 1)

	service.AddReceivers(Direct("some-user2@email.com"), Stream("stream-name", "topic-name"))
	assert.Len(service.receivers, 3)

	service.AddReceivers(Stream("another-stream-name", "topic-name"))
	assert.Len(service.receivers, 4)
}

func TestZulip_Send(t *testing.T) {
	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(nil, errors.New("some error"))

		service.client = mockClient

		service.AddReceivers(Direct("some-user@email.com"))

		err := service.Send(ctx, "subject", "message")
		assert.NotNil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("success direct message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockResponse := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id": 42,"msg": "","result": "success"}`)),
		}

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(mockResponse, nil)

		service.client = mockClient
		service.AddReceivers(Direct("some-user@email.com"))

		err := service.Send(ctx, "subject", "message")
		assert.Nil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("success stream message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockResponse := &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(`{"id": 42,"msg": "","result": "success"}`)),
		}

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(mockResponse, nil)

		service.client = mockClient
		service.AddReceivers(Stream("stream-name", "topic-name"))

		err := service.Send(ctx, "subject", "message")
		assert.Nil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("no receivers added", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		err := service.Send(ctx, "subject", "message")
		assert.Nil(err)
	})

	t.Run("non-exists receiver for a direct message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockResponse := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code": "BAD_REQUEST","msg": "Invalid email 'invalid@email.com'","result": "error"}`)),
		}

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(mockResponse, nil)

		service.client = mockClient
		service.AddReceivers(Direct("invalid@email.com"))

		err := service.Send(ctx, "subject", "message")
		assert.NotNil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("non-exists stream for a stream message", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockResponse := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code": "BAD_REQUEST","msg": "Invalid email 'invalid@email.com'","result": "error"}`)),
		}

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(mockResponse, nil)

		service.client = mockClient
		service.AddReceivers(Stream("invalid-stream-name", "topic-name"))

		err := service.Send(ctx, "subject", "message")
		assert.NotNil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("invalid response from API server", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		mockResponse := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code": "INTERNAL_SERVER_ERROR","msg": "Something went wrong'","result": "error"}`)),
		}

		mockClient := newMockZulipClient(t)
		mockClient.
			On("Message", mock.AnythingOfType("gozulipbot.Message")).
			Return(mockResponse, nil)

		service.client = mockClient
		service.AddReceivers(Direct("some-user@email.com"))

		err := service.Send(ctx, "subject", "message")
		assert.NotNil(err)
		mockClient.AssertExpectations(t)
	})

	t.Run("deadline exceeded", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		assert := require.New(t)

		service := New("your-org", "ZULIP_API_KEY", "email-bot@your-org.zulipchat.com")
		assert.NotNil(service)

		deadline := time.Now().Add(-5 * time.Second)
		ctx, cancelCtx := context.WithDeadline(ctx, deadline)
		defer cancelCtx()

		mockClient := newMockZulipClient(t)

		service.client = mockClient
		service.AddReceivers(Direct("some-user@email.com"))

		err := service.Send(ctx, "subject", "message")
		assert.NotNil(err)
	})
}
