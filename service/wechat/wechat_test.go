package wechat

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"github.com/stretchr/testify/require"
)

func TestAddReceivers(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		userIDs: []string{},
	}
	userIDs := []string{"User1ID", "User2ID", "User3ID"}
	svc.AddReceivers(userIDs...)

	assert.Equal(svc.userIDs, userIDs)
}

func TestSend(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		userIDs: []string{},
	}

	// test wechat message manager returning error
	mockMsgManager := new(mockWechatMessageManager)
	mockMsgManager.On("Send", message.NewCustomerTextMessage("UserID1", "subject\nmessage")).
		Return(errors.New("some error"))
	svc.messageManager = mockMsgManager
	svc.AddReceivers("UserID1")
	ctx := context.Background()
	err := svc.Send(ctx, "subject", "message")
	assert.NotNil(err)
	mockMsgManager.AssertExpectations(t)

	// test success and multiple receivers
	mockMsgManager = new(mockWechatMessageManager)
	mockMsgManager.On("Send", message.NewCustomerTextMessage("UserID1", "subject\nmessage")).
		Return(nil)
	mockMsgManager.On("Send", message.NewCustomerTextMessage("UserID2", "subject\nmessage")).
		Return(nil)
	svc.messageManager = mockMsgManager
	svc.AddReceivers("UserID1", "UserID2")
	err = svc.Send(ctx, "subject", "message")
	assert.Nil(err)
	mockMsgManager.AssertExpectations(t)
}
