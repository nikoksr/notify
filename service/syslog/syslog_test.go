package syslog

import (
	"context"
	"log/syslog"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestSyslog_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	// Test creating a local writer with invalid log priority.
	svc, err := New(-1, "")
	assert.Error(err)
	assert.Nil(svc)

	// Test creating a local writer successfully.
	svc, err = New(syslog.LOG_USER, "")
	assert.NoError(err)
	assert.NotNil(svc)
	err = svc.Close()
	assert.NoError(err)

	// Test creating a remote writer with invalid port.
	svc, err = NewFromDial("tcp", "localhost:99999", syslog.LOG_USER, "")
	assert.Error(err)
	assert.Nil(svc)

	// Test creating a remote writer successfully.
	svc, err = NewFromDial("", "", syslog.LOG_USER, "")
	assert.NoError(err)
	assert.NotNil(svc)
	err = svc.Close()
	assert.NoError(err)
}

func TestSyslog_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	ctx := context.Background()
	svc := new(Service)

	// Test syslog writer returning error
	mockClient := new(mockSyslogWriter)
	mockClient.
		On("Write", []byte("test: message")).
		Return(0, errors.New("some error"))
	svc.writer = mockClient

	err := svc.Send(ctx, "test", "message")
	assert.Error(err)
	mockClient.AssertExpectations(t)

	// Test syslog writer returning no error
	mockClient = new(mockSyslogWriter)
	mockClient.
		On("Write", []byte("test: message")).
		Return(0, nil)
	svc.writer = mockClient

	err = svc.Send(ctx, "test", "message")
	assert.NoError(err)
	mockClient.AssertExpectations(t)

	// Test closing the syslog writer with error
	mockClient = new(mockSyslogWriter)
	mockClient.
		On("Close").
		Return(errors.New("some error"))
	svc.writer = mockClient

	err = svc.Close()
	assert.Error(err)
	mockClient.AssertExpectations(t)

	// Test closing the syslog writer without error
	mockClient = new(mockSyslogWriter)
	mockClient.
		On("Close").
		Return(nil)
	svc.writer = mockClient

	err = svc.Close()
	assert.NoError(err)
	mockClient.AssertExpectations(t)
}
