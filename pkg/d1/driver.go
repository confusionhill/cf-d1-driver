package d1

import (
    "context"
    "database/sql/driver"
)

// DB is the abstraction boundary for the driver package.
// This keeps the implementation testable and allows mock-backed adapters.
type DB interface {
    ExecContext(ctx context.Context, query string, args ...any) (driver.Result, error)
    QueryContext(ctx context.Context, query string, args ...any) (driver.Rows, error)
    Batch(ctx context.Context, queries ...BatchQuery) ([]BatchResult, error)
    Ping(ctx context.Context) error
}

// QueryExecutor describes the subset of behavior needed for sqlx compatibility.
type QueryExecutor interface {
    ExecContext(ctx context.Context, query string, args ...any) (driver.Result, error)
    QueryContext(ctx context.Context, query string, args ...any) (driver.Rows, error)
}

// Driver is the exported driver type used by database/sql.
type Driver struct{}
