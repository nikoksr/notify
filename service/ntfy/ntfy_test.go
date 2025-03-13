package ntfy

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/AnthonyHewins/gotfy"
	"github.com/stretchr/testify/require"
)

func TestNtfy_New(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, _ := New(defaultServerURL, nil, nil)
	assert.NotNil(service)
}

func TestNtfy_NewWithPublishers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, _ := NewWithPublishers()
	assert.NotNil(service)

	pub1, _ := gotfy.NewPublisher(&url.URL{Host: "xyz.com"}, http.DefaultClient)
	pub2, _ := gotfy.NewPublisher(&url.URL{Host: "ahc.com"}, http.DefaultClient)
	service, _ = NewWithPublishers(*pub1, *pub2)
	assert.NotNil(service)
}

func TestNtfy_AddPublishers(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, _ := New(defaultServerURL, nil, nil)
	assert.NotNil(service)

	defaultServerU, _ := url.Parse(defaultServerURL)
	urlXyz, _ := url.Parse("https://xyz.com")
	urlAbc, _ := url.Parse("https://abc.com")

	pub1, _ := gotfy.NewPublisher(urlXyz, http.DefaultClient)
	pub2, _ := gotfy.NewPublisher(urlAbc, http.DefaultClient)
	defaultServerPub, _ := gotfy.NewPublisher(defaultServerU, http.DefaultClient)
	expected := []gotfy.Publisher{*defaultServerPub, *pub1, *pub2}

	service.AddPublishers(*pub1, *pub2)

	assert.Equal(expected, service.publishers)
}

func TestNtfy_send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	topic := "go-notify-test"
	service, _ := New(defaultServerURL, nil, &gotfy.Message{
		Topic: topic,
	})
	assert.NotNil(service)

	ctx := context.Background()

	remURL, _ := url.Parse("https://filesrv.lan/space.jpg")
	clickURL, _ := url.Parse("https://homecamera.lan/xasds1h2xsSsa/")

	content := gotfy.Message{
		Message:           "Disk space is low at 5.1 GB",
		Title:             "Low disk space alert",
		Tags:              []string{"warning", "cd"},
		Priority:          4,
		AttachURL:         remURL,
		AttachURLFilename: "diskspace.jpg",
		ClickURL:          clickURL,
	}

	defaultServer, _ := url.Parse(defaultServerURL)
	pub, _ := gotfy.NewPublisher(defaultServer, http.DefaultClient)

	service.AddMessage(content)
	err := service.send(ctx, pub, topic, "Disk space is low at 5.1 GB")
	assert.Error(err)

	err = service.send(ctx, nil, topic, "Disk space is low at 5.1 GB")
	assert.Error(err)

	pub, _ = gotfy.NewPublisher(&url.URL{Host: "xyz.com"}, http.DefaultClient)
	err = service.send(ctx, pub, topic, "sample_body")
	assert.Error(err)

	err = service.send(ctx, pub, topic, "Disk space is low at 5.1 GB")
	assert.Error(err)
}

func TestNtfy_Send(t *testing.T) {
	t.Parallel()

	assert := require.New(t)

	service, _ := New("", nil, &gotfy.Message{
		Topic:   "go-notify",
		Message: "Disk space is low at 5.1 GB",
	})
	assert.NotNil(service)

	ctx := context.Background()
	topic := "go-notify"
	title := "Low disk space alert"

	service.AddMessage(gotfy.Message{
		Topic: topic,
	})

	err := service.Send(ctx, title, "Disk space is low at 5.1 GB")
	assert.NoError(err)

	service.AddMessage(gotfy.Message{
		Topic: "",
	})

	err = service.Send(ctx, title, "Disk space is low at 5.1 GB")
	assert.Error(err)
}
