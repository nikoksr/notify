package msgerrors

import (
	"errors"
	"strings"
)

// MsgError wraps a slice of errors and provides convenient
// syntax for processing individual errors without breaking
// the service
type MsgError struct {
	err []string
}

// New instantiates a MsgError with a nil Err field
func New() *MsgError {
	return &MsgError{err: make([]string, 0, 0)}
}

// Append safely adds a new error to our list
func (m *MsgError) Append(err error) {
	if err != nil {
		m.err = append(m.err, err.Error())
	}
}

// Error provides the string output of our error
// This allows MsgError to implement the core "error"
// interface
func (m *MsgError) Error() string {
	if len(m.err) == 0 {
		return ""
	}
	return strings.Join(m.err, "\n")
}

// Errors provides conventient syntax for returning
// an error containing all information about the
// accumulated errors in code. It will return
// nil if no errors have been collated
func (m *MsgError) Errors() error {
	if len(m.err) == 0 {
		return nil
	}
	return errors.New(m.Error())
}
