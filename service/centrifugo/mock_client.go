package centrifugo

import (
	"context"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

type MockClient struct {
	PublishFunc func(ctx context.Context, channel string, data []byte) (centrifuge.PublishResult, error)
	Closed      bool
}

func (m *MockClient) Publish(ctx context.Context, channel string, data []byte) (centrifuge.PublishResult, error) {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, channel, data)
	}
	return centrifuge.PublishResult{}, nil
}

func (m *MockClient) Close() {
	m.Closed = true
}
