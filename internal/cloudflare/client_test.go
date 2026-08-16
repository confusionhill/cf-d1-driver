package cloudflare_test

import (
	"context"
	"testing"

	cf "github.com/confusionhill/cf-d1-driver/internal/cloudflare"
)

func TestNewClient(t *testing.T) {
	client := cf.NewClient("acc_123", "db_456", "token_123", "https://api.cloudflare.com/client/v4")
	if client == nil {
		t.Fatal("expected client, got nil")
	}
}

func TestClientQueryRequest(t *testing.T) {
	client := cf.NewClient("acc_123", "db_456", "token_123", "https://api.cloudflare.com/client/v4")

	req := cf.QueryRequest{SQL: "SELECT 1", Params: []any{1}}
	if req.SQL == "" {
		t.Fatal("expected sql to be set")
	}

	if _, err := client.Query(context.Background(), req); err == nil {
		t.Fatal("expected query to fail for uninitialized transport")
	}
}

func TestClientBatchRequest(t *testing.T) {
	client := cf.NewClient("acc_123", "db_456", "token_123", "https://api.cloudflare.com/client/v4")

	req := cf.BatchRequest{Batch: []cf.QueryRequest{{SQL: "INSERT INTO test (value) VALUES (?)", Params: []any{"ok"}}}}
	if _, err := client.Batch(context.Background(), req); err == nil {
		t.Fatal("expected batch to fail for uninitialized transport")
	}
}
