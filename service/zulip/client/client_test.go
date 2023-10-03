package client

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSend(t *testing.T) {
	t.Parallel()

	t.Run("send=success", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			email, apiKey, _ := req.BasicAuth()
			assert.Equal(t, email, "email", "invalid email provided")
			assert.Equal(t, apiKey, "apiKey", "invalid apiKey provided")

			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
			    "id":42,
			    "msg": "",
			    "result": "success"
			}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithCreds("email", "apiKey"),
			WithTimeout(10*time.Second),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			Type:    "stream",
			To:      "test",
			Topic:   "general",
			Content: "some content goes here",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if resp.Result != "success" {
			t.Fatalf("invalid response: %v", err)
		}
	})

	t.Run("send=failure", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			email, apiKey, _ := req.BasicAuth()
			assert.Equal(t, email, "email", "invalid email provided")
			assert.Equal(t, apiKey, "apiKey", "invalid apiKey provided")

			rw.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithCreds("email", "apiKey"),
			WithTimeout(10*time.Second),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		resp, err := client.Send(&Message{
			Type:    "stream",
			To:      "test",
			Topic:   "general",
			Content: "some content goes here",
		})
		if err == nil {
			t.Fatal("expected error but got nil")
		}
		if resp != nil {
			t.Fatalf("expected nil response\ngot: %v response", resp)
		}
	})

	t.Run("send=invalid_token", func(t *testing.T) {
		t.Parallel()

		_, err := NewClient(
			WithCreds("", ""),
			WithTimeout(10*time.Second),
		)
		if err == nil {
			t.Fatal("expected error but got nil")
		}
	})

	t.Run("send=invalid_message", func(t *testing.T) {
		t.Parallel()

		client, err := NewClient(
			WithBaseURL(""),
			WithCreds("email", "apiKey"),
			WithTimeout(10*time.Second),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		_, err = client.Send(&Message{})
		if err == nil {
			t.Fatal("expected error but go nil")
		}
	})
}

func TestSendWithContext(t *testing.T) {
	t.Run("send_context=success", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			email, apiKey, _ := req.BasicAuth()
			assert.Equal(t, email, "email", "invalid email provided")
			assert.Equal(t, apiKey, "apiKey", "invalid apiKey provided")

			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
			    "id":42,
			    "msg": "",
			    "result": "success"
			}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithCreds("email", "apiKey"),
			WithTimeout(10*time.Second),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ctx := context.Background()
		resp, err := client.SendWithContext(ctx, &Message{
			Type:    "stream",
			To:      "test",
			Topic:   "general",
			Content: "some content goes here",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp.Result != "success" {
			t.Fatalf("invalid response: %v", err)
		}
	})

	t.Run("send_context=timeout", func(t *testing.T) {
		t.Parallel()

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			email, apiKey, _ := req.BasicAuth()
			assert.Equal(t, email, "email", "invalid email provided")
			assert.Equal(t, apiKey, "apiKey", "invalid apiKey provided")

			time.Sleep(time.Millisecond * 100)

			rw.WriteHeader(http.StatusOK)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, `{
			    "id":42,
			    "msg": "",
			    "result": "success"
			}`)
		}))
		defer server.Close()

		client, err := NewClient(
			WithBaseURL(server.URL),
			WithCreds("email", "apiKey"),
			WithTimeout(100*time.Millisecond),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*50)
		defer cancel()
		_, err = client.SendWithContext(ctx, &Message{
			Type:    "stream",
			To:      "test",
			Topic:   "general",
			Content: "some content goes here",
		})
		if err == nil {
			t.Fatalf("no context timeout")
		}
	})
}
