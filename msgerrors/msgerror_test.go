package msgerrors

import (
	"errors"
	"strings"
	"testing"
)

func TestMsgError(t *testing.T) {
	msgErr := New()

	if msgErr.Error() != "" {
		t.Errorf("Expected empty string. Recieved %s", msgErr.Error())
	}
	if msgErr.Errors() != nil {
		t.Errorf("Expected nil. Recieved %s", msgErr.Errors())
	}

	err1 := errors.New("Error 1")
	err2 := errors.New("Error 2")
	msgErr.Append(err1)
	msgErr.Append(err2)

	if len(msgErr.err) != 2 {
		t.Errorf("Expected length 2. Recieved %d", len(msgErr.err))
	}

	concatenatedErr := strings.Join([]string{"Error 1", "Error 2"}, "\n")
	if msgErr.Error() != concatenatedErr {
		t.Errorf("Expected %s. \n Recieved %s", concatenatedErr, msgErr.Error())
	}

	if msgErr.Errors() == nil {
		t.Errorf("Expected %s. \n Recieved %s", errors.New(concatenatedErr), msgErr.Errors())
	}
}
