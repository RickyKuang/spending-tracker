// Package apiclient is a thin HTTP client for the spending-tracker backend API.
package apiclient

import (
	"context"
	"net/http"
)

// Client talks to the spending-tracker backend API over HTTP.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates a Client for the given backend base URL and bearer token.
func New(baseURL, token string) *Client {
	return &Client{
		baseURL:    baseURL,
		token:      token,
		httpClient: &http.Client{},
	}
}

// do sends a JSON request to path and decodes the JSON response into out.
// body may be nil for requests without a body; out may be nil to discard the response body.
func (c *Client) do(ctx context.Context, method, path string, body, out any) error {
	// TODO: implement
	return nil
}
