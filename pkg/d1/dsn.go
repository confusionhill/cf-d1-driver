package d1

import (
	"fmt"
	"net/url"
	"strings"
)

// Config contains the parsed Cloudflare D1 connection configuration.
type Config struct {
	AccountID  string
	DatabaseID string
	APIToken   string
	APIBaseURL string
}

// ParseDSN parses a DSN in the form:
// d1://account_id/database_id?token=...&api_base=...
func ParseDSN(raw string) (Config, error) {
	if raw == "" {
		return Config{}, fmt.Errorf("d1: empty DSN")
	}

	u, err := url.Parse(raw)
	if err != nil {
		return Config{}, fmt.Errorf("d1: invalid DSN: %w", err)
	}

	if u.Scheme != "d1" {
		return Config{}, fmt.Errorf("d1: unsupported scheme %q", u.Scheme)
	}

	accountID := strings.Trim(u.Host, "/")
	databaseID := strings.TrimPrefix(u.Path, "/")

	if accountID == "" {
		return Config{}, fmt.Errorf("d1: missing account_id in DSN")
	}

	if databaseID == "" {
		return Config{}, fmt.Errorf("d1: missing database_id in DSN")
	}

	cfg := Config{
		AccountID:  accountID,
		DatabaseID: databaseID,
		APIToken:   u.Query().Get("token"),
		APIBaseURL: u.Query().Get("api_base"),
	}

	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = "https://api.cloudflare.com/client/v4"
	}

	if cfg.APIToken == "" {
		return Config{}, fmt.Errorf("d1: missing token in DSN")
	}

	return cfg, nil
}
