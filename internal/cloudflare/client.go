package cloudflare

import (
	"context"
	"fmt"
	"net/http"
)

// QueryRequest models a single D1 SQL query payload.
type QueryRequest struct {
	SQL    string `json:"sql"`
	Params []any  `json:"params,omitempty"`
}

// BatchRequest models a D1 atomic batch payload.
type BatchRequest struct {
	Batch []QueryRequest `json:"batch"`
}

// Client is the abstraction boundary for the Cloudflare HTTP implementation.
type Client interface {
	Query(ctx context.Context, req QueryRequest) error
	Batch(ctx context.Context, req BatchRequest) error
}

// client is the concrete Cloudflare implementation used behind the interface.
type client struct {
	accountID  string
	databaseID string
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Cloudflare client.
func NewClient(accountID, databaseID, token, baseURL string) Client {
	return &client{
		accountID:  accountID,
		databaseID: databaseID,
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// Query executes a single D1 query request.
func (c *client) Query(ctx context.Context, req QueryRequest) error {
	if c == nil || c.httpClient == nil {
		return fmt.Errorf("cloudflare: client is not initialized")
	}
	if req.SQL == "" {
		return fmt.Errorf("cloudflare: SQL query cannot be empty")
	}
	return fmt.Errorf("cloudflare: query not implemented in this scaffold")
}

// Batch executes a D1 batch request.
func (c *client) Batch(ctx context.Context, req BatchRequest) error {
	if c == nil || c.httpClient == nil {
		return fmt.Errorf("cloudflare: client is not initialized")
	}
	if len(req.Batch) == 0 {
		return fmt.Errorf("cloudflare: batch query cannot be empty")
	}
	return fmt.Errorf("cloudflare: batch not implemented in this scaffold")
}
