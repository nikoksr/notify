package client

import (
	"errors"
	"reflect"
)

type Message struct {
	Type    string
	To      any
	Topic   string
	Content string
}

// Validate returns an error if the message is not well-formed.
func (m *Message) Validate() error {
	switch m.Type {
	case "direct":
	case "stream":
	case "private":
		break
	default:
		return errors.New("invalid message type. Type should be 'direct','stream' or 'private'")
	}

	switch reflect.TypeOf(m.To) {
	case reflect.TypeOf(""):
	case reflect.TypeOf(1):
	case reflect.SliceOf(reflect.TypeOf("")):
	case reflect.SliceOf(reflect.TypeOf(1)):
		break
	default:
		return errors.New("invalid message to. To should be string, int, []string or []int")
	}

	return nil
}
