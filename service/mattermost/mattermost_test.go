package mattermost

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

const url = "https://host.mattermost.com"

func TestService_New(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	service := New(url)
	assert.NotNil(service)
}

func TestService_LoginWithCredentials(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	service := New(url)
	// Test responses
	mockClient := newMockHttpClient(t)
	mockClient.
		On("Send", context.TODO(), "fake-loginID", "fake-password").Return(nil)
	service.loginClient = mockClient
	// test call
	err := service.LoginWithCredentials(context.TODO(), "fake-loginID", "fake-password")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test errors
	// Test responses
	mockClient = newMockHttpClient(t)
	mockClient.
		On("Send", context.TODO(), "fake-loginID", "").Return(errors.New("empty password"))
	service.loginClient = mockClient
	// test call
	err = service.LoginWithCredentials(context.TODO(), "fake-loginID", "")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)
}

func TestService_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New(url)
	assert.NotNil(service)

	service.AddReceivers("yfgstwuisnshydhd")
	assert.Equal(1, len(service.channelIDs))

	service.AddReceivers("yfgstwuisnshydhd", "nwudneyfrwqjs")
	assert.Equal(3, len(service.channelIDs))

	hooks := []string{"yfgstwuisnshydhd", "nwudneyfrwqjs"}
	service.channelIDs = []string{}
	service.AddReceivers(hooks...)
	assert.Equal(service.channelIDs, hooks)
}

func TestService_Send(t *testing.T) {
	t.Parallel()
	assert := require.New(t)

	service := New(url)
	channelID := "yfgstwuisnshydhd"
	service.channelIDs = append(service.channelIDs, channelID)

	// Test responses
	mockClient := newMockHttpClient(t)
	mockClient.
		On("Send", context.TODO(), channelID, "fake-sub\nfake-msg").Return(nil)
	service.messageClient = mockClient
	// test call
	err := service.Send(context.TODO(), "fake-sub", "fake-msg")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test error on Send
	// Test responses
	mockClient = newMockHttpClient(t)
	mockClient.
		On("Send", context.TODO(), channelID, "fake-sub\nfake-msg").Return(errors.New("internal error"))
	service.messageClient = mockClient
	// test call
	err = service.Send(context.TODO(), "fake-sub", "fake-msg")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)
}
