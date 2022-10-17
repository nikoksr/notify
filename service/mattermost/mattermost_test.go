package mattermost

import (
	"context"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/stretchr/testify/require"
)

const url = "https://host.mattermost.com"

func TestMatrix_New(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	service := New(url)
	assert.NotNil(service)
}

func TestService_LoginWithCredentials(t *testing.T) {
	t.Parallel()
	assert := require.New(t)
	// Test response
	mockClient := newMockMattermostClient(t)
	mockClient.
		On("Login", "fake-loginID", "fake-password").Return(
		&model.User{}, &model.Response{Error: nil},
	)
	service := New(url)
	service.client = mockClient
	err := service.LoginWithCredentials("fake-loginID", "fake-password")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test error on Send
	mockClient = newMockMattermostClient(t)
	mockClient.
		On("Login", "fake-loginID", "").Return(
		&model.User{}, &model.Response{
			Error: &model.AppError{
				StatusCode: 401,
				Message:    "missing password",
			},
		})
	service = New(url)
	service.client = mockClient
	err = service.LoginWithCredentials("fake-loginID", "")
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

	ctx := context.Background()
	channelID := "yfgstwuisnshydhd"

	service := New(url)
	err := service.Send(ctx, "fake-subject", "fake-message")
	assert.Nil(err)

	// Test response
	mockClient := newMockMattermostClient(t)
	validPost := &model.Post{
		ChannelId: channelID,
		Message:   "fake-subject\nfake-message",
	}
	mockClient.
		On("CreatePost", validPost).Return(validPost,
		&model.Response{
			Error: nil,
		})
	service = New(url)
	service.client = mockClient
	service.AddReceivers(channelID)
	err = service.Send(ctx, "fake-subject", "fake-message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// Test error on Send
	mockClient = newMockMattermostClient(t)
	inValidPost := &model.Post{
		ChannelId: channelID,
		Message:   "\n",
	}
	mockClient.
		On("CreatePost", inValidPost).Return(inValidPost,
		&model.Response{
			Error: &model.AppError{
				StatusCode: 500,
				Message:    "internal error",
			},
		})
	service = New(url)
	service.client = mockClient
	service.AddReceivers(channelID)
	err = service.Send(ctx, "", "")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)
}
