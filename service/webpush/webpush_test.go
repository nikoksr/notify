package webpush

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/google/go-cmp/cmp"
)

// Allows us to simulate an error returned from the server on a per-request basis.
const headerTestError = "X-Test-Error"

// checkFunc is a function that checks a request and returns an error if the check fails.
type checkFunc func(r *http.Request) error

// checkMethod returns a checkFunc that checks the request method.
func checkMethod(method string) func(r *http.Request) error {
	return func(r *http.Request) error {
		if r.Method != method {
			return fmt.Errorf("unexpected method: %s", r.Method)
		}
		return nil
	}
}

// checkHeader returns a checkFunc that checks the request header.
func checkHeader(key, value string) func(r *http.Request) error {
	return func(r *http.Request) error {
		if r.Header.Get(key) != value {
			return fmt.Errorf("unexpected %s header: %s", key, r.Header.Get(key))
		}
		return nil
	}
}

// defaultChecks is the default set of checks used for testing.
func defaultChecks() []checkFunc {
	return []checkFunc{
		checkMethod("POST"),
		checkHeader("Content-Type", "application/octet-stream"),
		checkHeader("Content-Encoding", "aes128gcm"),
	}
}

// newWebpushHandlerWithChecks returns a new http.Handler that checks the request against the given checks and returns
// a 400 if any of them fail.
func newWebpushHandlerWithChecks(checks ...checkFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, check := range checks {
			if err := check(r); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
		}

		// This allows us to simulate an error returned from the server on a per-request basis
		if r.Header.Get(headerTestError) != "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(r.Header.Get(headerTestError)))
			return
		}

		w.WriteHeader(http.StatusCreated)
	})
}

//nolint:gochecknoglobals // These are used across tests and read-only.
var vapidPublicKey, vapidPrivateKey string

// TestMain sets up a test server to handle the requests.
func TestMain(m *testing.M) {
	// Generate a VAPID key pair
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		slog.Error("failed to generate VAPID keys", "error", err) //nolint:sloglint // Test setup, context not available
		os.Exit(1)
	}

	vapidPublicKey = publicKey
	vapidPrivateKey = privateKey

	os.Exit(m.Run())
}

func getValidSubscription() webpush.Subscription {
	return webpush.Subscription{
		Keys: webpush.Keys{
			P256dh: "BNNL5ZaTfK81qhXOx23-wewhigUeFb632jN6LvRWCFH1ubQr77FE_9qV1FuojuRmHP42zmf34rXgW80OvUVDgTk",
			Auth:   "zqbxT6JKstKSY9JKibZLSQ",
		},
	}
}

func getInvalidSubscription() Subscription {
	return Subscription{
		Keys: webpush.Keys{
			P256dh: "AAA",
			Auth:   "BBB",
		},
	}
}

//nolint:tparallel // We need to run the sub-tests sequentially.
func TestService_Send(t *testing.T) {
	t.Parallel()

	type fields struct {
		subscriptions   []webpush.Subscription
		vapidPublicKey  string
		vapidPrivateKey string
	}
	type args struct {
		subject string
		message string
		options webpush.Options // Bind those to the context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		handler http.Handler
		wantErr bool
	}{
		{
			name: "Send a message with options",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getValidSubscription(),
				},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{
					TTL:     60,           // Should set the TTL header
					Topic:   "test-topic", // Should set the Topic header
					Urgency: UrgencyHigh,  // Should set the Urgency header
				},
			},
			handler: newWebpushHandlerWithChecks(
				append(
					defaultChecks(),
					checkHeader("TTL", "60"),
					checkHeader("Topic", "test-topic"),
					checkHeader("Urgency", string(UrgencyHigh)),
				)...,
			),
			wantErr: false,
		},
		{
			name: "Send a message with no options",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getValidSubscription(),
				},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: false,
		},
		{
			name: "Send a message with no options and no subscriptions",
			fields: fields{
				subscriptions:   []webpush.Subscription{},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: false,
		},
		{
			name: "Send a message with no options and no subscriptions",
			fields: fields{
				subscriptions:   []webpush.Subscription{},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: false,
		},
		{
			name: "Send a message with no vapid public key",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getValidSubscription(),
				},
				vapidPublicKey:  "",
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: false, // Yes, does not cause an error
		},
		{
			name: "Send a message with no vapid private key",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getValidSubscription(),
				},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: "",
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: true, // Empty VAPID private key causes parse error
		},
		{
			name: "Send a message with no vapid keys",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getValidSubscription(),
				},
				vapidPublicKey:  "",
				vapidPrivateKey: "",
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: true, // Empty VAPID keys cause parse error
		},
		{
			name: "Send a message with invalid subscription",
			fields: fields{
				subscriptions: []webpush.Subscription{
					getInvalidSubscription(),
				},
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				subject: "subject",
				message: "message",
				options: webpush.Options{},
			},
			handler: newWebpushHandlerWithChecks(defaultChecks()...),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeWebpushServer := httptest.NewServer(tt.handler)
			defer fakeWebpushServer.Close()

			s := New(tt.fields.vapidPublicKey, tt.fields.vapidPrivateKey)

			for _, subscription := range tt.fields.subscriptions {
				subscription.Endpoint = fakeWebpushServer.URL
				s.AddReceivers(subscription)
			}

			ctx := WithOptions(context.Background(), tt.args.options)
			err := s.Send(ctx, tt.args.subject, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_withOptions(t *testing.T) {
	t.Parallel()

	type fields struct {
		vapidPublicKey  string
		vapidPrivateKey string
	}
	type args struct {
		options Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Options
	}{
		{
			name: "with options but no VAPID keys",
			fields: fields{
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				options: Options{
					RecordSize:      4096,
					Subscriber:      "test-subscriber",
					Topic:           "test-topic",
					TTL:             100,
					Urgency:         UrgencyHigh,
					VAPIDPublicKey:  "", // should be ignored
					VAPIDPrivateKey: "", // should be ignored
				},
			},
			want: Options{
				RecordSize:      4096,
				Subscriber:      "test-subscriber",
				Topic:           "test-topic",
				TTL:             100,
				Urgency:         UrgencyHigh,
				VAPIDPublicKey:  vapidPublicKey,
				VAPIDPrivateKey: vapidPrivateKey,
			},
		},
		{
			name: "with options and VAPID keys",
			fields: fields{
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			args: args{
				options: Options{
					RecordSize:      4096,
					Subscriber:      "test-subscriber",
					Topic:           "test-topic",
					TTL:             100,
					Urgency:         UrgencyHigh,
					VAPIDPublicKey:  "test-public-key",  // should be used
					VAPIDPrivateKey: "test-private-key", // should be used
				},
			},
			want: Options{
				RecordSize:      4096,
				Subscriber:      "test-subscriber",
				Topic:           "test-topic",
				TTL:             100,
				Urgency:         UrgencyHigh,
				VAPIDPublicKey:  "test-public-key",
				VAPIDPrivateKey: "test-private-key",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := New(tt.fields.vapidPublicKey, tt.fields.vapidPrivateKey)

			if got := s.withOptions(tt.args.options); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("withOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		vapidPublicKey  string
		vapidPrivateKey string
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "New with VAPID keys",
			args: args{
				vapidPublicKey:  vapidPublicKey,
				vapidPrivateKey: vapidPrivateKey,
			},
			want: &Service{
				subscriptions: []webpush.Subscription{},
				options: Options{
					VAPIDPublicKey:  vapidPublicKey,
					VAPIDPrivateKey: vapidPrivateKey,
				},
			},
		},
		{
			name: "New without VAPID keys",
			args: args{
				vapidPublicKey:  "",
				vapidPrivateKey: "",
			},
			want: &Service{
				subscriptions: []webpush.Subscription{},
				options: Options{
					VAPIDPublicKey:  "",
					VAPIDPrivateKey: "",
				},
			},
		},
	}

	opts := []cmp.Option{
		cmp.AllowUnexported(Service{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := New(tt.args.vapidPublicKey, tt.args.vapidPrivateKey)

			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("New() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestService_AddReceivers(t *testing.T) {
	t.Parallel()

	type fields struct {
		subscriptions []webpush.Subscription
	}
	type args struct {
		subscriptions []Subscription
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "AddReceivers",
			fields: fields{
				subscriptions: []webpush.Subscription{},
			},
			args: args{
				subscriptions: []Subscription{
					{
						Endpoint: "https://fcm.googleapis.com/fcm/send/dQw4w9WgXcQ",
						Keys: webpush.Keys{
							Auth:   "auth",
							P256dh: "p256dh",
						},
					},
				},
			},
		},
		{
			name: "AddReceivers with multiple subscriptions",
			fields: fields{
				subscriptions: []webpush.Subscription{},
			},
			args: args{
				subscriptions: []Subscription{
					{
						Endpoint: "https://example.com/push1",
						Keys: webpush.Keys{
							Auth:   "auth",
							P256dh: "p256dh",
						},
					},
					{
						Endpoint: "https://example.com/push2",
						Keys: webpush.Keys{
							Auth:   "auth",
							P256dh: "p256dh",
						},
					},
					{
						Endpoint: "https://example.com/push3",
						Keys: webpush.Keys{
							Auth:   "auth",
							P256dh: "p256dh",
						},
					},
				},
			},
		},
	}

	opts := []cmp.Option{
		cmp.AllowUnexported(Service{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := New("", "")
			s.AddReceivers(tt.args.subscriptions...)

			if diff := cmp.Diff(s.subscriptions, tt.args.subscriptions, opts...); diff != "" {
				t.Errorf("AddReceivers() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_ContextBinding(t *testing.T) {
	t.Parallel()

	type fields struct {
		ctx context.Context
	}
	type args struct {
		data    map[string]any
		options Options
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Bind data",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				data: map[string]any{
					"title": "Test",
				},
			},
		},
		{
			name: "Bind options",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				options: Options{
					Topic:   "test",
					Urgency: UrgencyHigh,
					TTL:     60,
				},
			},
		},
		{
			name: "Bind data and options",
			fields: fields{
				ctx: context.Background(),
			},
			args: args{
				data: map[string]any{
					"title": "Test",
				},
				options: Options{
					Topic:   "test",
					Urgency: UrgencyHigh,
					TTL:     60,
				},
			},
		},
		{
			name: "Bind nothing", // Make sure nothing panics
			fields: fields{
				ctx: context.Background(),
			},
			args: args{},
		},
	}

	opts := []cmp.Option{
		cmp.AllowUnexported(Service{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCtx := WithData(tt.fields.ctx, tt.args.data)
			gotCtx = WithOptions(gotCtx, tt.args.options)

			if gotCtx == nil {
				t.Error("gotCtx is nil")
			}

			gotData := dataFromContext(gotCtx)
			gotOptions := optionsFromContext(gotCtx)

			if diff := cmp.Diff(gotData, tt.args.data, opts...); diff != "" {
				t.Errorf("WithData() mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(gotOptions, tt.args.options, opts...); diff != "" {
				t.Errorf("WithOptions() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_payloadFromContext(t *testing.T) {
	t.Parallel()

	type args struct {
		subject string
		message string
		data    map[string]any
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Payload with only subject and message",
			args: args{
				subject: "test",
				message: "test",
			},
			want: []byte(`{"subject":"test","message":"test"}`),
		},
		{
			name: "Payload with subject, message and data",
			args: args{
				subject: "test",
				message: "test",
				data: map[string]any{
					"title": "Test",
				},
			},
			want: []byte(`{"subject":"test","message":"test","data":{"title":"Test"}}`),
		},
	}

	opts := []cmp.Option{
		cmp.AllowUnexported(Service{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			if tt.args.data != nil {
				ctx = WithData(ctx, tt.args.data)
			}

			got, err := payloadFromContext(ctx, tt.args.subject, tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("payloadFromContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(got, tt.want, opts...); diff != "" {
				t.Errorf("payloadFromContext() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
