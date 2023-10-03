package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL contains example base URL of zulip service.
	DefaultBaseURL = "https://yourZulipDomain.zulipchat.com"

	// DefaultTimeout duration in second
	DefaultTimeout time.Duration = 30 * time.Second
)

// ErrInvalidCreds occurs if API key or Email is not set.
var ErrInvalidCreds = errors.New("client credentials are invalid")

// Client abstracts the interaction between the application server and the
// Zulip server via HTTP protocol. The developer must obtain an API key from the
// Zulip's personal settings page and pass it to the `Client` so that it can
// perform authorized requests on the application server's behalf.
// To send a message to one or more devices use the Client's Send.
//
// If the `HTTP` field is nil, a zeroed http.Client will be allocated and used
// to send messages.
type Client struct {
	email   string
	apiKey  string
	baseURL string
	client  *http.Client
	timeout time.Duration
}

type Option func(*Client) error

// NewClient creates new Zulip Client based on opts passed and
// with default endpoint and http client.
func NewClient(opt Option) (*Client, error) {
	c := &Client{
		baseURL: DefaultBaseURL,
		client:  &http.Client{},
		timeout: DefaultTimeout,
	}
	if err := opt(c); err != nil {
		return nil, err
	}

	if c.apiKey == "" || c.email == "" {
		return nil, ErrInvalidCreds
	}

	return c, nil
}

type Response struct {
	ID     int    `json:"id"`
	Msg    string `json:"msg"`
	Result string `json:"result"`
	Code   string `json:"code"`
}

// SendWithContext sends a message to the Zulip server without retrying in case of service
// unavailability. A non-nil error is returned if a non-recoverable error
// occurs (i.e. if the response status is not "200 OK").
// Behaves just like regular send, but uses external context.
func (c *Client) SendWithContext(ctx context.Context, msg *Message) (*Response, error) {
	// validate
	if err := msg.Validate(); err != nil {
		return nil, err
	}

	return c.send(ctx, msg)
}

// Send sends a message to the Zulip server without retrying in case of service
// unavailability. A non-nil error is returned if a non-recoverable error
// occurs (i.e. if the response status is not "200 OK").
func (c *Client) Send(msg *Message) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.SendWithContext(ctx, msg)
}

// send sends a request.
func (c *Client) send(ctx context.Context, msg *Message) (*Response, error) {
	// set the message data
	data := url.Values{}
	data.Set("type", msg.Type)
	data.Set("to", fmt.Sprintf("%v", msg.To))
	data.Set("topic", msg.Topic)
	data.Set("content", msg.Content)

	// create request
	url, _ := url.JoinPath(c.baseURL, "/api/v1/messages")
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	// add headers
	req.SetBasicAuth(c.email, c.apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// execute request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check response status
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode >= http.StatusInternalServerError {
			return nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
		}
		return nil, fmt.Errorf("%d error: %s", resp.StatusCode, resp.Status)
	}

	// build return
	response := new(Response)
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return nil, err
	}

	return response, nil
}
