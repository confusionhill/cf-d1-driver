package d1

import (
	"errors"
	"fmt"
	"strings"
)

// BatchQuery models a single statement entry in a D1 atomic batch.
type BatchQuery struct {
	SQL    string
	Params []any
}

// Validate ensures the batch query is meaningful.
func (q BatchQuery) Validate() error {
	if strings.TrimSpace(q.SQL) == "" {
		return errors.New("d1: batch query SQL cannot be empty")
	}
	return nil
}

// BatchResult contains the result for a single statement in a batch.
type BatchResult struct {
	Meta    QueryMeta        `json:"meta,omitempty"`
	Results []map[string]any `json:"results,omitempty"`
	Success bool             `json:"success,omitempty"`
}

// QueryMeta contains metadata returned by Cloudflare D1 for each query.
type QueryMeta struct {
	ChangedDB       bool         `json:"changed_db,omitempty"`
	Changes         int64        `json:"changes,omitempty"`
	Duration        float64      `json:"duration,omitempty"`
	LastRowID       int64        `json:"last_row_id,omitempty"`
	RowsRead        int64        `json:"rows_read,omitempty"`
	RowsWritten     int64        `json:"rows_written,omitempty"`
	ServedByClo     string       `json:"served_by_colo,omitempty"`
	ServedByPrimary bool         `json:"served_by_primary,omitempty"`
	ServedByRegion  string       `json:"served_by_region,omitempty"`
	SizeAfter       int64        `json:"size_after,omitempty"`
	Timings         QueryTimings `json:"timings,omitempty"`
}

// QueryTimings is a subset of the D1 metadata response.
type QueryTimings struct {
	SQLDurationMS int64 `json:"sql_duration_ms,omitempty"`
}

// ValidateBatch validates a full batch payload.
func ValidateBatch(queries []BatchQuery) error {
	if len(queries) == 0 {
		return fmt.Errorf("d1: batch cannot be empty")
	}

	for _, q := range queries {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}
