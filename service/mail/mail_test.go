package mail

import (
	"crypto/tls"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMail_newEmailHtml(t *testing.T) {
	t.Parallel()

	text := "test"
	m := New("foo", "server")
	email := m.newEmail("test", text)

	assert.False(t, m.usePlainText)
	assert.Equal(t, []byte(nil), email.Text)
	assert.Equal(t, []byte(text), email.HTML)
}

func TestMail_newEmailText(t *testing.T) {
	t.Parallel()

	text := "test"
	m := New("foo", "server")
	m.BodyFormat(PlainText)
	email := m.newEmail("test", text)

	assert.True(t, m.usePlainText)
	assert.Equal(t, []byte(text), email.Text)
	assert.Equal(t, []byte(nil), email.HTML)
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
	assert.Nil(t, m.smtpAuth)

	m.AuthenticateSMTP("test", "test", "test", "test")
	assert.NotNil(t, m.smtpAuth)
}

func TestMail_SetTLS(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")
	assert.False(t, m.useTLS)
	assert.Nil(t, m.tlsConfig)

	m.SetTLS(&tls.Config{ServerName: "server"})
	assert.True(t, m.useTLS)
	assert.NotNil(t, m.tlsConfig)
}

func TestMail_UnSetTLS(t *testing.T) {
	t.Parallel()

	m := New("foo", "server")
	assert.False(t, m.useTLS)
	assert.Nil(t, m.tlsConfig)

	m.SetTLS(&tls.Config{ServerName: "server"})
	assert.True(t, m.useTLS)
	assert.NotNil(t, m.tlsConfig)

	m.UnSetTLS()
	assert.False(t, m.useTLS)
}
