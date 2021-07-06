package twilio

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAddReceivers(t *testing.T) {
	assert := require.New(t)

	svc := &Service{
		contacts: []string{},
	}
	contacts := []string{"Contact1", "Contact2", "Contact3"}
	svc.AddReceivers(contacts...)

	assert.Equal(svc.contacts, contacts)
}
