package zulip

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	z "github.com/nikoksr/notify/service/zulip/client"
)

func TestZulip_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("server-base-url", "bot-email", "api-key")
	assert.NotNil(service)
	assert.Nil(err)
}

func TestDirectHook(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	receiver := Direct("test@gmail.com")

	assert.Equal(receiver._type, "direct")
	assert.Equal(receiver._to, "test@gmail.com")
	assert.Equal(receiver._topic, "<ignored>")

	assert.NotNil(receiver)
	assert.IsType(new(Receiver), receiver)
}

func TestStreamHook(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	receiver := Stream("group", "topic")

	assert.Equal(receiver._type, "stream")
	assert.Equal(receiver._to, "group")
	assert.Equal(receiver._topic, "topic")

	assert.NotNil(receiver)
	assert.IsType(new(Receiver), receiver)
}

func TestZulip_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, _ := New("server-base-url", "bot-email", "api-key")

	service.AddReceivers(Direct("test@gmail.com"))
	assert.Len(service.recv, 1)

	service.AddReceivers(Direct("test2@gmail.com"), Stream("stream", "topic"))
	assert.Len(service.recv, 3)

	service.AddReceivers(Stream("stream2", "topic2"))
	assert.Len(service.recv, 4)
}

func TestZulip_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	t.Run("send_message=direct", func(t *testing.T) {
		t.Parallel()

		service := &Service{}

		mockZulipClient := newMockZulipClient(t)
		mockZulipClient.On("Send", &z.Message{
			Type:    "direct",
			To:      "test@gmail.com",
			Topic:   "<ignored>",
			Content: "subject message",
		}).Return(&z.Response{
			ID:     0,
			Msg:    "",
			Result: "",
			Code:   "",
		}, nil)

		service.client = mockZulipClient

		service.AddReceivers(Direct("test@gmail.com"))
		err := service.Send(context.Background(), "subject", "message")
		assert.Nil(err)
		mockZulipClient.AssertExpectations(t)
	})

	t.Run("send_message=stream", func(t *testing.T) {
		t.Parallel()

		service := &Service{}

		mockZulipClient := newMockZulipClient(t)
		mockZulipClient.On("Send", &z.Message{
			Type:    "stream",
			To:      "group",
			Topic:   "topic",
			Content: "subject message",
		}).Return(&z.Response{
			ID:     0,
			Msg:    "",
			Result: "",
			Code:   "",
		}, nil)

		service.client = mockZulipClient

		service.AddReceivers(Stream("group", "topic"))
		err := service.Send(context.Background(), "subject", "message")
		assert.Nil(err)
		mockZulipClient.AssertExpectations(t)
	})
}
