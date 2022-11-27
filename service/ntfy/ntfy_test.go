package ntfy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNtfy_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)
}

func TestNtfy_NewWithServers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := NewWithServers()
	assert.NotNil(service)

	service = NewWithServers("xyz.com", "abc.com")
	assert.NotNil(service)
}

func TestNtfy_AddReceivers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	rec := []string{"https://rec1.sh/", "https://rec2.sh/", "https://rec3.sh/"}
	expected := []string{"https://ntfy.sh/", "https://rec1.sh/", "https://rec2.sh/", "https://rec3.sh/"}
	service.AddReceivers(rec...)

	assert.Equal(service.serverURLs, expected)
}

func TestNtfy_send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	service.client = defaultHTTPClient()

	ctx := context.Background()
	serverURL := DefaultServerURL
	topic := "pushkar"
	content := `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"tags": ["warning","cd"],
		"priority": 4,
		"attach": "https://filesrv.lan/space.jpg",
		"filename": "diskspace.jpg",
		"click": "https://homecamera.lan/xasds1h2xsSsa/",
		"actions": [{ "action": "view", "label": "Admin panel", "url": "https://filesrv.lan/admin" }]
	}`

	err := service.send(ctx, serverURL, topic, content)
	assert.Nil(err)

	err = service.send(ctx, "", topic, content)
	assert.NotNil(err)

	err = service.send(ctx, "xyz", topic, content)
	assert.NotNil(err)

	err = service.send(ctx, serverURL, topic, "sample_body")
	assert.NotNil(err)

	content = `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"priority": 4222,
	}`

	err = service.send(ctx, serverURL, topic, content)
	assert.NotNil(err)

}

func TestNtfy_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service := New()
	assert.NotNil(service)

	ctx := context.Background()
	topic := "pushkar"
	content := `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"tags": ["warning","cd"],
		"priority": 4,
		"attach": "https://filesrv.lan/space.jpg",
		"filename": "diskspace.jpg",
		"click": "https://homecamera.lan/xasds1h2xsSsa/",
		"actions": [{ "action": "view", "label": "Admin panel", "url": "https://filesrv.lan/admin" }]
	}`

	err := service.Send(ctx, topic, content)
	assert.Nil(err)

	content = `{
		"message": "Disk space is low at 5.1 GB",
		"title": "Low disk space alert",
		"priority": 4222,
	}`

	err = service.Send(ctx, topic, content)
	assert.NotNil(err)
}
