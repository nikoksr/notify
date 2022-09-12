package lark

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLark_NewCustomAppService(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := NewCustomAppService("", "")
	assert.NotNil(service)
}

func TestLark_AddReceivers(t *testing.T) {
	t.Parallel()

	xs := []*ReceiverID{
		OpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
		UserID("8335aga2"),
	}
	svc := NewCustomAppService("", "")
	svc.AddReceivers(xs...)

	assert.ElementsMatch(t, svc.receiveIDs, xs)

	// Test if adding more receivers afterwards works.
	ys := []*ReceiverID{
		UnionID("on_cad4860e7af114fb4ff6c5d496d1dd76"),
		Email("xyz@example.com"),
		ChatID("oc_a0553eda9014c201e6969b478895c230"),
	}
	svc.AddReceivers(ys...)

	assert.ElementsMatch(t, svc.receiveIDs, append(xs, ys...))
}

func TestLark_SendCustomApp(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	assert := assert.New(t)

	tests := []*ReceiverID{
		OpenID("ou_c99c5f35d542efc7ee492afe11af19ef"),
		UserID("8335aga2"),
		UnionID("on_cad4860e7af114fb4ff6c5d496d1dd76"),
		Email("xyz@example.com"),
		ChatID("oc_a0553eda9014c201e6969b478895c230"),
	}

	// First, test for when the sender returns an error.
	for _, tt := range tests {
		mockSendToer := newMockSendToer(t)
		mockSendToer.
			On("SendTo", "subject", "message", tt.id, string(tt.typ)).
			Return(errors.New(""))

		svc := NewCustomAppService("", "")
		svc.cli = mockSendToer

		svc.AddReceivers(tt)
		err := svc.Send(ctx, "subject", "message")
		assert.NotNil(err)

		mockSendToer.AssertExpectations(t)
	}

	// Then test for when the sender does not return an error.
	for _, tt := range tests {
		mockSendToer := newMockSendToer(t)
		mockSendToer.
			On("SendTo", "subject", "message", tt.id, string(tt.typ)).
			Return(nil)

		svc := NewCustomAppService("", "")
		svc.cli = mockSendToer

		svc.AddReceivers(tt)
		err := svc.Send(ctx, "subject", "message")
		assert.Nil(err)

		mockSendToer.AssertExpectations(t)
	}
}
