package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

// QueryResult is the D1 result payload for a single statement.
type QueryResult struct {
	Meta    QueryMeta        `json:"meta,omitempty"`
	Results []map[string]any `json:"results,omitempty"`
	Success bool             `json:"success,omitempty"`
}

// QueryMeta is the metadata returned by D1.
type QueryMeta struct {
	ChangedDB     bool   `json:"changed_db,omitempty"`
	Changes       int64  `json:"changes,omitempty"`
	Duration      float64 `json:"duration,omitempty"`
	LastRowID     int64  `json:"last_row_id,omitempty"`
	RowsRead      int64  `json:"rows_read,omitempty"`
	RowsWritten   int64  `json:"rows_written,omitempty"`
	SizeAfter     int64  `json:"size_after,omitempty"`
	ServedByClo   string `json:"served_by_colo,omitempty"`
	ServedByPrimary bool `json:"served_by_primary,omitempty"`
	ServedByRegion string `json:"served_by_region,omitempty"`
}

// Client is the abstraction boundary for the Cloudflare HTTP implementation.
type Client interface {
	Query(ctx context.Context, req QueryRequest) ([]QueryResult, error)
	Batch(ctx context.Context, req BatchRequest) ([]QueryResult, error)
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

func (c *client) endpoint() string {
	base := c.baseURL
	if base == "" {
		base = "https://api.cloudflare.com/client/v4"
	}
	return fmt.Sprintf("%s/accounts/%s/d1/database/%s/query", strings.TrimRight(base, "/"), url.PathEscape(c.accountID), url.PathEscape(c.databaseID))
}

func (c *client) do(ctx context.Context, payload any) ([]QueryResult, error) {
	if c == nil || c.httpClient == nil {
		return nil, fmt.Errorf("cloudflare: client is not initialized")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cloudflare: marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("cloudflare: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cloudflare: request failed: %w", err)
	}
	defer resp.Body.Close()

	payloadBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cloudflare: read response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("cloudflare: http %d: %s", resp.StatusCode, string(payloadBytes))
	}

	var apiResp struct {
		Success bool          `json:"success"`
		Result  []QueryResult `json:"result"`
		Errors  []APIError    `json:"errors"`
	}
	if err := json.Unmarshal(payloadBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("cloudflare: decode response: %w", err)
	}

	if !apiResp.Success {
		if len(apiResp.Errors) > 0 {
			return nil, fmt.Errorf("cloudflare: %s", apiResp.Errors[0].Message)
		}
		return nil, fmt.Errorf("cloudflare: request failed")
	}

	return apiResp.Result, nil
}

// Query executes a single D1 query request.
func (c *client) Query(ctx context.Context, req QueryRequest) ([]QueryResult, error) {
	if req.SQL == "" {
		return nil, fmt.Errorf("cloudflare: SQL query cannot be empty")
	}
	return c.do(ctx, req)
}

// Batch executes a D1 batch request.
func (c *client) Batch(ctx context.Context, req BatchRequest) ([]QueryResult, error) {
	if len(req.Batch) == 0 {
		return nil, fmt.Errorf("cloudflare: batch query cannot be empty")
	}
	return c.do(ctx, req)
}
