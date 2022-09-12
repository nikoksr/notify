package fcm

import (
	"context"
	"errors"
	"testing"

	"github.com/appleboy/go-fcm"
	"github.com/stretchr/testify/require"
)

func TestFCM_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, err := New("server-api-key")
	assert.NotNil(service)
	assert.Nil(err)
}

func TestFCM_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc := &Service{
		deviceTokens: []string{},
	}
	deviceTokens := []string{"Token1", "Token2", "Token3"}
	svc.AddReceivers(deviceTokens...)

	assert.Equal(svc.deviceTokens, deviceTokens)
}

func TestFCM_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	svc := &Service{
		deviceTokens: []string{},
	}

	mockToken := "deviceToken"
	mockData := map[string]interface{}{
		"foo": "bar",
	}

	// test fcm client send
	mockClient := newMockFcmClient(t)
	mockClient.On("SendWithRetry", &fcm.Message{
		To: mockToken,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 0).Return(&fcm.Response{Success: 1}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockToken)
	ctx := context.Background()
	err := svc.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// test fcm client send with data
	mockClient = newMockFcmClient(t)
	mockClient.On("SendWithRetry", &fcm.Message{
		To:   mockToken,
		Data: mockData,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 0).Return(&fcm.Response{Success: 1}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockToken)
	ctx = context.Background()
	ctxWithData := context.WithValue(ctx, DataKey, mockData)
	err = svc.Send(ctxWithData, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// test fcm client send with data and retries
	mockClient = newMockFcmClient(t)
	mockClient.On("SendWithRetry", &fcm.Message{
		To:   mockToken,
		Data: mockData,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 3).Return(&fcm.Response{Success: 1}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockToken)
	ctx = context.Background()
	ctxWithData = context.WithValue(ctx, DataKey, mockData)
	ctxWithDataAndRetries := context.WithValue(ctxWithData, RetriesKey, 3)
	err = svc.Send(ctxWithDataAndRetries, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)

	// test fcm client returning error
	mockClient = newMockFcmClient(t)
	mockClient.On("SendWithRetry", &fcm.Message{
		To:   mockToken,
		Data: mockData,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 3).Return(nil, errors.New("some error"))
	svc.client = mockClient
	svc.AddReceivers(mockToken)
	ctx = context.Background()
	ctxWithData = context.WithValue(ctx, DataKey, mockData)
	ctxWithDataAndRetries = context.WithValue(ctxWithData, RetriesKey, 3)
	err = svc.Send(ctxWithDataAndRetries, "subject", "message")
	assert.NotNil(err)
	mockClient.AssertExpectations(t)

	// test fcm client multiple receivers
	anotherMockToken := "another_device_token"
	mockClient = newMockFcmClient(t)
	mockClient.On("SendWithRetry", &fcm.Message{
		To:   mockToken,
		Data: mockData,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 3).Return(&fcm.Response{Success: 1}, nil)
	mockClient.On("SendWithRetry", &fcm.Message{
		To:   anotherMockToken,
		Data: mockData,
		Notification: &fcm.Notification{
			Title: "subject",
			Body:  "message",
		},
	}, 3).Return(&fcm.Response{Success: 1}, nil)
	svc.client = mockClient
	svc.AddReceivers(mockToken, anotherMockToken)
	ctx = context.Background()
	ctxWithData = context.WithValue(ctx, DataKey, mockData)
	ctxWithDataAndRetries = context.WithValue(ctxWithData, RetriesKey, 3)
	err = svc.Send(ctxWithDataAndRetries, "subject", "message")
	assert.Nil(err)
	mockClient.AssertExpectations(t)
}

func TestFCM_GetMessageData(t *testing.T) {
	t.Parallel()

	assert := require.New(t)
	mockData := map[string]interface{}{
		"foo": "bar",
	}

	// test without data
	ctx := context.Background()
	data, ok := getMessageData(ctx)
	assert.Nil(data)
	assert.False(ok)

	// test with invalid type of data
	ctx = context.Background()
	ctxWithInvalidTypeOfData := context.WithValue(ctx, DataKey, "invalid_type_of_data")
	data, ok = getMessageData(ctxWithInvalidTypeOfData)
	assert.Nil(data)
	assert.False(ok)

	// test with invalid data key
	ctx = context.Background()
	invalidDataKey := struct {
		Key string
	}{
		Key: "invalid_data_key",
	}
	ctxWithInvalidDataKey := context.WithValue(ctx, invalidDataKey, mockData)
	data, ok = getMessageData(ctxWithInvalidDataKey)
	assert.Nil(data)
	assert.False(ok)

	// test with data
	ctx = context.Background()
	ctxWithData := context.WithValue(ctx, DataKey, mockData)
	data, ok = getMessageData(ctxWithData)
	assert.Equal(mockData, data)
	assert.True(ok)
}

func TestFCM_GetMessageRetryAttempts(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	// test without retries
	ctx := context.Background()
	retryAttempts := getMessageRetryAttempts(ctx)
	assert.Equal(0, retryAttempts)

	// test with invalid type of retries
	ctx = context.Background()
	ctxWithInvalidTypeOfRetries := context.WithValue(ctx, RetriesKey, "invalid_type_of_retries")
	retryAttempts = getMessageRetryAttempts(ctxWithInvalidTypeOfRetries)
	assert.Equal(0, retryAttempts)

	// test with invalid retries key
	ctx = context.Background()
	invalidRetriesKey := struct {
		Key string
	}{
		Key: "invalid_retries_key",
	}
	ctxWithInvalidRetriesKey := context.WithValue(ctx, invalidRetriesKey, 3)
	retryAttempts = getMessageRetryAttempts(ctxWithInvalidRetriesKey)
	assert.Equal(0, retryAttempts)

	// test with retries
	ctx = context.Background()
	ctxWithRetries := context.WithValue(ctx, RetriesKey, 3)
	retryAttempts = getMessageRetryAttempts(ctxWithRetries)
	assert.Equal(3, retryAttempts)
}
