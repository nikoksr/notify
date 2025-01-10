package pagerduty_test

import (
	"fmt"
	"testing"

	gopagerduty "github.com/PagerDuty/go-pagerduty"
	"github.com/nikoksr/notify/service/pagerduty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_NewConfig(t *testing.T) {
	t.Parallel()

	config := pagerduty.NewConfig()

	want := &pagerduty.Config{
		Receivers:        []string{},
		NotificationType: "incident",
	}

	require.Equal(t, want, config)
}

func TestConfig_OK(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  *pagerduty.Config
		wantErr error
	}{
		{
			name: "ok_basic_config",
			config: &pagerduty.Config{
				FromAddress:      "sender@domain.com",
				Receivers:        []string{"AB1234", "CD5678"},
				NotificationType: "incident",
			},
			wantErr: nil,
		},
		{
			name: "ok_complete_config",
			config: &pagerduty.Config{
				FromAddress:      "sender@domain.com",
				Receivers:        []string{"AB1234", "CD5678"},
				Urgency:          "high",
				PriorityID:       "P1234",
				NotificationType: "incident",
			},
			wantErr: nil,
		},
		{
			name: "missing_from_address",
			config: &pagerduty.Config{
				Receivers: []string{"AB1234", "CD5678"},
			},
			wantErr: fmt.Errorf("from address is required"),
		},
		{
			name: "invalid_from_address",
			config: &pagerduty.Config{
				FromAddress: "senderdomain.com",
				Receivers:   []string{"AB1234", "CD5678"},
			},
			wantErr: fmt.Errorf("from address is invalid: mail: missing '@' or angle-addr"),
		},
		{
			name: "missing_receivers",
			config: &pagerduty.Config{
				FromAddress: "sender@domain.com",
			},
			wantErr: fmt.Errorf("at least one receiver is required"),
		},
		{
			name: "missing_notification_type",
			config: &pagerduty.Config{
				FromAddress: "sender@domain.com",
				Receivers:   []string{"AB1234", "CD5678"},
			},
			wantErr: fmt.Errorf("notification type is required"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			err := test.config.OK()
			if test.wantErr == nil {
				require.NoError(t, err)
				return
			}

			require.EqualError(t, err, test.wantErr.Error())
		})
	}
}

func TestConfig_Priority(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		config *pagerduty.Config
		want   *gopagerduty.APIReference
	}{
		{
			name:   "default_priority",
			config: &pagerduty.Config{},
			want:   nil,
		},
		{
			name: "priority",
			config: &pagerduty.Config{
				PriorityID: "P1234",
			},
			want: &gopagerduty.APIReference{
				ID:   "P1234",
				Type: pagerduty.APIPriorityReference,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := test.config.PriorityReference()
			require.Equal(t, test.want, got)
		})
	}
}

func TestConfig_SetFromAddress(t *testing.T) {
	t.Parallel()

	config := pagerduty.NewConfig()
	require.NotNil(t, config)

	assert.Empty(t, config.FromAddress)

	config.SetFromAddress("sender@domain.com")
	assert.Equal(t, "sender@domain.com", config.FromAddress)
}

func TestConfig_AddReceivers(t *testing.T) {
	t.Parallel()

	config := pagerduty.NewConfig()
	require.NotNil(t, config)

	assert.Empty(t, config.Receivers)

	config.AddReceivers("AB1234", "CD5678")
	assert.Equal(t, []string{"AB1234", "CD5678"}, config.Receivers)
}

func TestConfig_SetNotificationType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		config           *pagerduty.Config
		notificationType string
		want             string
	}{
		{
			name:             "empty_notification_type",
			config:           &pagerduty.Config{},
			notificationType: "",
			want:             "incident",
		},
		{
			name:             "set_notification_type",
			config:           &pagerduty.Config{},
			notificationType: "notification",
			want:             "notification",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			test.config.SetNotificationType(test.notificationType)
			assert.Equal(t, test.want, test.config.NotificationType)
		})
	}
}

func TestConfig_SetPriorityID(t *testing.T) {
	t.Parallel()

	config := pagerduty.NewConfig()
	require.NotNil(t, config)

	assert.Empty(t, config.PriorityID)

	config.SetPriorityID("P1234")
	assert.Equal(t, "P1234", config.PriorityID)
}

func TestConfig_SetUrgency(t *testing.T) {
	t.Parallel()

	config := pagerduty.NewConfig()
	require.NotNil(t, config)

	assert.Empty(t, config.Urgency)

	config.SetUrgency("high")
	assert.Equal(t, "high", config.Urgency)
}
