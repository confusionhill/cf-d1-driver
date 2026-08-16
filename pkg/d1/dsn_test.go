package d1

import "testing"

func TestParseDSN(t *testing.T) {
	cfg, err := ParseDSN("d1://acc_123/db_456?token=secret&api_base=https://api.cloudflare.com/client/v4")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.AccountID != "acc_123" {
		t.Fatalf("expected account id, got %q", cfg.AccountID)
	}

	if cfg.DatabaseID != "db_456" {
		t.Fatalf("expected database id, got %q", cfg.DatabaseID)
	}

	if cfg.APIToken != "secret" {
		t.Fatalf("expected api token, got %q", cfg.APIToken)
	}
}
