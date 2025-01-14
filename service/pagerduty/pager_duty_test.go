package pagerduty_test

import (
	"context"
	"errors"
	"testing"

	gopagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/nikoksr/notify/service/pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockClient struct {
	mock.Mock
}

func (m *mockClient) CreateIncidentWithContext(
	ctx context.Context,
	from string,
	options *gopagerduty.CreateIncidentOptions,
) (*gopagerduty.Incident, error) {
	args := m.Called(ctx, from, options)

	if err := args.Error(1); err != nil {
		return nil, err
	}

	incident, isIncident := args.Get(0).(*gopagerduty.Incident)
	if !isIncident {
		return nil, errors.New("unexpected type for first argument")
	}

	return incident, nil
}

func TestPagerDuty_New(t *testing.T) {
	t.Parallel()

	t.Run("successful_new", func(t *testing.T) {
		t.Parallel()

		service, err := pagerduty.New("fake_token")
		require.NoError(t, err)

		want := &pagerduty.PagerDuty{
			Config: pagerduty.NewConfig(),
			Client: gopagerduty.NewClient("fake_token"),
		}

		assert.Equal(t, want, service)
	})

	t.Run("fail_new_invalid_token", func(t *testing.T) {
		t.Parallel()

		_, err := pagerduty.New("")
		require.EqualError(t, err, "access token is required")
	})
}

func TestPagerDuty_Send(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		receivers     []string
		subject       string
		message       string
		mockSetup     func(m *mockClient)
		expectedError string
		expectedCall  bool // whether the mock should be called
	}{
		{
			name:         "successful_send_to_multiple_receivers",
			receivers:    []string{"AB1234", "CD5678"},
			subject:      "Test Subject",
			message:      "Test Message",
			expectedCall: true,
			mockSetup: func(m *mockClient) {
				m.On("CreateIncidentWithContext", mock.Anything, "sender@domain.com", mock.Anything).
					Return(&gopagerduty.Incident{}, nil)
			},
		},
		{
			name:         "successful_send_to_single_receivers",
			receivers:    []string{"AB1234"},
			subject:      "Test Subject",
			message:      "Test Message",
			expectedCall: true,
			mockSetup: func(m *mockClient) {
				m.On("CreateIncidentWithContext", mock.Anything, "sender@domain.com", mock.Anything).
					Return(&gopagerduty.Incident{}, nil)
			},
		},
		{
			name:         "unsuccessful_send",
			receivers:    []string{"AB1234"},
			subject:      "Test Subject",
			message:      "Test Message",
			expectedCall: true,
			mockSetup: func(m *mockClient) {
				m.On("CreateIncidentWithContext", mock.Anything, "sender@domain.com", mock.Anything).
					Return(nil, errors.New("failed to create incident"))
			},
			expectedError: "create pager duty incident: failed to create incident",
		},
		{
			name:    "fail_send_no_receivers",
			subject: "Test Subject",
			message: "Test Message",
			mockSetup: func(m *mockClient) {
				m.On("CreateIncidentWithContext", mock.Anything, "sender@domain.com", mock.Anything).
					Return(nil, errors.New("invalid configuration: at least one receiver is required"))
			},
			expectedError: "invalid configuration: at least one receiver is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockClient := new(mockClient)

			service, err := pagerduty.New("fake_token")
			require.NoError(t, err)
			require.NotNil(t, service)

			service.AddReceivers(test.receivers...)
			service.SetFromAddress("sender@domain.com")

			test.mockSetup(mockClient)

			service.Client = mockClient

			err = service.Send(context.Background(), test.subject, test.message)

			if test.expectedError != "" {
				require.EqualError(t, err, test.expectedError)
			} else {
				require.NoError(t, err)
			}

			if test.expectedCall {
				mockClient.AssertExpectations(t)
			}
		})
	}
}
