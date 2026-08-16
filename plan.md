# Cloudflare D1 Go Driver API Specification

## Goal

Provide a Go `database/sql` driver for Cloudflare D1, compatible with `sqlx`, while respecting D1's actual remote API semantics.

This specification is intentionally strict: it reflects how Cloudflare D1 behaves, not how a local SQLite driver behaves.

## Core design principle

D1 does not support SQL transaction commands such as `BEGIN`, `COMMIT`, or `ROLLBACK`.

The D1 API exposes a batch execution model instead:

```json
{
  "batch": [
    { "sql": "INSERT ...", "params": [...] },
    { "sql": "UPDATE ...", "params": [...] }
  ]
}
```

Therefore:

- Single statement execution maps to `/query` with `sql` + `params`
- Multi-statement atomic execution maps to `/query` with `batch`
- `database/sql` transactions are not implemented as SQL transactions
- a separate batch helper is required for atomic multi-step operations

## Driver package layout

```text
cf-d1-go-driver/
  d1/
    driver.go
    conn.go
    stmt.go
    rows.go
    result.go
    dsn.go
    batch.go
    errors.go
  internal/
    client/
      client.go
      types.go
      query.go
      response.go
```

## Public API

### Registration

```go
import "database/sql"

func init() {
    sql.Register("d1", &d1.Driver{})
}
```

### DSN

```go
// Examples:
//   d1://acc_123/db_456?token=cf_token
//   d1://cf_token@acc_123/db_456
//   d1://acc_123/db_456?token=cf_token&api_base=https://api.cloudflare.com/client/v4
```

### Package-level helper

```go
func Open(dsn string) (*sql.DB, error)
```

This is a convenience wrapper around `sql.Open("d1", dsn)`.

## Driver types

```go
package d1

type Driver struct{}

type Config struct {
    AccountID string
    DatabaseID string
    APIToken string
    APIBaseURL string
    HTTPClient *http.Client
}

type Conn struct {
    cfg    Config
    client *Client
}

type Stmt struct {
    conn *Conn
    query string
}

type Rows struct {
    columns []string
    rows    []map[string]any
    index   int
}

type Result struct {
    RowsAffected int64
    LastInsertID int64
}
```

## Batch API

This is the D1-native transaction-style feature.

```go
package d1

type BatchQuery struct {
    SQL    string
    Params []any
}

type BatchResult struct {
    Meta    QueryMeta
    Results []map[string]any
    Success bool
}

func (c *Conn) Batch(ctx context.Context, queries ...BatchQuery) ([]BatchResult, error)
```

### Rules

- `Batch` sends a request with `batch` field to D1
- if the whole batch fails, the method returns an error
- the query list is atomic at the D1 API level
- `Batch` is the supported replacement for SQL transaction behavior

## Query request model

### Single query request

```go
type QueryRequest struct {
    SQL    string `json:"sql"`
    Params []any  `json:"params,omitempty"`
}
```

### Batch query request

```go
type BatchRequest struct {
    Batch []QueryRequest `json:"batch"`
}
```

## Response model

```go
type D1Response struct {
    Result  []QueryResult `json:"result"`
    Errors  []APIError    `json:"errors"`
    Messages []APIMessage `json:"messages"`
    Success bool          `json:"success"`
}

type QueryResult struct {
    Meta    QueryMeta        `json:"meta,omitempty"`
    Results []map[string]any `json:"results,omitempty"`
    Success bool             `json:"success,omitempty"`
}

type QueryMeta struct {
    ChangedDB   bool    `json:"changed_db,omitempty"`
    Changes     int64   `json:"changes,omitempty"`
    Duration    float64 `json:"duration,omitempty"`
    LastRowID   int64   `json:"last_row_id,omitempty"`
    RowsRead    int64   `json:"rows_read,omitempty"`
    RowsWritten int64   `json:"rows_written,omitempty"`
    ServedByClo string  `json:"served_by_colo,omitempty"`
    ServedByPrimary bool `json:"served_by_primary,omitempty"`
    ServedByRegion string `json:"served_by_region,omitempty"`
    SizeAfter   int64   `json:"size_after,omitempty"`
    Timings     QueryTimings `json:"timings,omitempty"`
}

type QueryTimings struct {
    SQLDurationMS int64 `json:"sql_duration_ms,omitempty"`
}

type APIError struct {
    Code int    `json:"code"`
    Message string `json:"message"`
    DocumentationURL string `json:"documentation_url,omitempty"`
    Source SourceInfo `json:"source,omitempty"`
}

type SourceInfo struct {
    Pointer string `json:"pointer,omitempty"`
}
```

## `database/sql` compatibility expectations

### Supported

```go
func (c *Conn) Prepare(query string) (driver.Stmt, error)
func (c *Conn) Close() error
func (c *Conn) Ping(ctx context.Context) error
func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error)
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error)
```

### Explicitly unsupported

```go
func (c *Conn) Begin() (driver.Tx, error)
```

Behavior:

```go
return nil, errors.New("d1: SQL transactions are unsupported; use Batch() for atomic multi-statement operations")
```

This is required because D1 does not accept SQL transaction commands.

## Argument conversion

The driver must convert Go values into JSON-safe values that D1 accepts.

### Supported value conversions

- `string`
- `[]byte`
- `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `time.Time` => ISO-8601 string
- `nil` => `null`
- pointers to supported primitives => dereference before encoding

### Unsupported

- custom structs as values in params
- channel types
- function types
- unsupported interface values

The driver should return a clear conversion error in these cases.

## Execution semantics

### Single-statement `ExecContext`

- build `QueryRequest{SQL: query, Params: convertedArgs}`
- call D1 `/query`
- read the first `QueryResult`
- set `RowsAffected` from `meta.changes`
- set `LastInsertID` from `meta.last_row_id`
- return `driver.Result`

### Single-statement `QueryContext`

- build `QueryRequest{SQL: query, Params: convertedArgs}`
- call D1 `/query`
- read the first `QueryResult`
- convert `results` array into `driver.Rows`
- ensure column names are taken from each row object or the result schema if available

### Batch execution

- build `BatchRequest{Batch: []QueryRequest{...}}`
- call D1 `/query`
- decode all `result` entries
- return slice of `BatchResult`

## SQL compatibility considerations

The driver supports SQL syntax that D1 supports, but not all of SQLite's SQL-level transaction features.

### Supported SQL features

- `SELECT`
- `INSERT`
- `UPDATE`
- `DELETE`
- `CREATE TABLE`
- `DROP TABLE`
- `ALTER TABLE` if D1 supports it at the time of use
- parameterized queries using `?`

### Not supported by the driver contract

- `BEGIN TRANSACTION`
- `COMMIT`
- `ROLLBACK`
- driver-level `Tx` behavior modeled after SQLite

## Error semantics

The driver should wrap all D1 API failures with a typed Go error:

```go
type Error struct {
    Code    int
    Message string
    DocURL  string
    Source  string
}

func (e *Error) Error() string
```

This is emitted for:

- HTTP 4xx / 5xx failures
- JSON `success: false`
- malformed query responses
- unsupported feature usage

## sqlx behavior

The driver should be used as:

```go
import (
    "github.com/jmoiron/sqlx"
)

func main() {
    db, err := sqlx.Open("d1", "account_id=acc_123&database_id=db_456&token=cf_token")
    if err != nil {
        panic(err)
    }
    defer db.Close()
}
```

### `sqlx` expectations

- `Get`, `Select`, `Exec`, and `Queryx` should work normally for single statements
- `db.Beginx()` is not supported and should return a clear error
- D1 batch operations should be used for atomic multi-step writes

## Example usage

```go
package main

import (
    "context"
    "fmt"

    d1 "github.com/your-org/cf-d1-go-driver/d1"
)

func main() {
    db, err := sql.Open("d1", "account_id=acc_123&database_id=db_456&token=cf_token")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    _, err = db.ExecContext(context.Background(), "INSERT INTO users (name) VALUES (?)", "alice")
    if err != nil {
        panic(err)
    }

    _, err = db.QueryContext(context.Background(), "SELECT * FROM users WHERE name = ?", "alice")
    if err != nil {
        panic(err)
    }

    _, err = d1.NewConn(...).Batch(context.Background(),
        d1.BatchQuery{SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"bob"}},
        d1.BatchQuery{SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"bobby", "bob"}},
    )
    if err != nil {
        panic(err)
    }

    fmt.Println("ok")
}
```

## Implementation priorities

### Phase 1

- DSN parsing
- HTTP client for D1 query endpoint
- `ExecContext` and `QueryContext`
- error translation

### Phase 2

- `Rows` conversion
- `sqlx` validation
- `Ping`

### Phase 3

- `Batch` API support
- transaction unsupported documentation
- examples and README

## Final compliance statement

This driver must respect the Cloudflare D1 API contract:

- SQL statements are submitted via `/query`
- atomic multi-step logic uses `/query` with `batch`
- SQL `BEGIN/COMMIT/ROLLBACK` are not supported by D1 and therefore not emulated by the driver
- the driver exposes the correct abstraction layer without pretending SQLite transaction semantics exist when they do not

5. `sqlx` example validation
6. only then consider batching and advanced behavior

This keeps the first version focused and avoids trying to emulate SQLite features that D1 does not expose over HTTP.
