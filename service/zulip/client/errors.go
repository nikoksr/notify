package client

import "errors"

// ErrInvalidCreds occurs if API key or Email is not set.
var ErrInvalidCreds = errors.New("client credentials are invalid")

// ErrInvalidMessageType occurs if Message Type is not 'direct','stream' or 'private'
var ErrInvalidMessageType = errors.New("invalid message type. Type should be 'direct','stream' or 'private'")

// ErrInvalidMessageTo occurs if Message To is not string, int, []string, []int
var ErrInvalidMessageTo = errors.New("invalid message to. To should be string, int, []string or []int")

// ErrInvalidMessageTopic occurs if Message Topic is empty
var ErrInvalidMessageTopic = errors.New("invalid message topic. Topic should be a non-empty string")

// ErrInvalidMessageContent occurs if Message Content is empty
var ErrInvalidMessageContent = errors.New("invalid message content. Content should be a non-empty string")
