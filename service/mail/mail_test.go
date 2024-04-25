package mail

import (
	"testing"

	"github.com/alecthomas/assert"
)

func TestMail_newEmailHtml(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")

	assert.False(t, m.usePlainText)
}

func TestMail_newEmailText(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")
	m.BodyFormat(PlainText)

	assert.True(t, m.usePlainText)

}

func TestMail_AddReceivers(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")
	m.AddReceivers("test")

	assert.Len(t, m.receiverAddresses, 1)
	assert.Equal(t, "test", m.receiverAddresses[0])
}

func TestMail_AuthenticateSMTP(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")

	m.AddAuthentication("test", "test")
	assert.NotNil(t, m.pass)
	assert.NotNil(t, m.user)
}

func TestMail_AddHeaders(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")

	m.AddHeader("test", "test")
	assert.Len(t, m.headers, 1)
	assert.Equal(t, "test", m.headers.Get("test"))
}
