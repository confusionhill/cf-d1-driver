# cf-d1-go-driver

> Unofficial project. This driver is still experimental and may change without notice.

A Cloudflare D1 Go driver built for `database/sql` compatibility and `sqlx` usage.

## Design goals

- Standard Go project layout
- Clean interfaces and dependency injection
- Testable packages
- D1-native atomic batch semantics instead of fake SQL transactions
- TDD-friendly structure with mocks

## Install

```bash
go get github.com/confusionhill/cf-d1-driver@v0.0.1
```

## Add to your project

The package is imported directly in your app. For a driver-only `database/sql` setup, use a blank import so the driver registers itself:

```go
import (
    "database/sql"
    "fmt"

    _ "github.com/confusionhill/cf-d1-driver/pkg/d1"
)

func main() {
    db, err := sql.Open("d1", "d1://account_id/database_id?token=secret")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    rows, err := db.Query("SELECT 1 AS ok")
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    for rows.Next() {
        var ok int
        if err := rows.Scan(&ok); err != nil {
            panic(err)
        }
        fmt.Println("result:", ok)
    }
}
```

For `sqlx`, use the separate example module in `examples/sqlx_app`, which binds the driver and imports the package for consumer-side integration:

```bash
cd examples/sqlx_app
go run .
```

## Supported and unsupported behavior

### Supported

- `database/sql` driver registration via `sql.Register("d1", &Driver{})`
- `sql.Open("d1", dsn)` and `db.Ping()` style usage
- simple query execution with `QueryContext` / `ExecContext`
- result scanning into Go values using the standard `database/sql` row API
- basic D1 HTTP request/response handling for query and batch operations
- `sqlx` consumer integration through `sqlx.BindDriver("d1", sqlx.QUESTION)`
- D1 batch semantics for atomic multi-statement writes

### Not supported or intentionally limited

- SQL transaction commands such as `BEGIN`, `COMMIT`, and `ROLLBACK`
- prepared statement support via `Prepare` / `Stmt` in the current driver layer
- full SQLite compatibility semantics outside the subset D1 exposes over its API
- streaming or advanced database features that Cloudflare D1 does not expose through its HTTP API
- any behavior that requires native SQLite engine features not provided by D1

This driver intentionally reflects Cloudflare D1 semantics instead of pretending to be a full SQLite driver.

## Official D1 API references

This implementation is based on the Cloudflare D1 HTTP API contract used for:

- database queries
- statement execution
- batch operations
- metadata returned by D1 results

Official docs:

- Cloudflare D1 docs: https://developers.cloudflare.com/d1/
- D1 API reference: https://developers.cloudflare.com/api/operations/cloudflare-d1-query-database
- D1 batch operations reference: https://developers.cloudflare.com/api/operations/cloudflare-d1-create-database

The driver follows the D1 API contract rather than SQLite semantics, so the implementation is intentionally aligned with what Cloudflare exposes.

## Transaction model

D1 does not support SQL `BEGIN`, `COMMIT`, or `ROLLBACK` statements.

Use the batch API for atomic multi-step writes:

```go
batch := []d1.BatchQuery{
    {SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}},
    {SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"alicia", "alice"}},
}
```

## Sample apps

### Driver-only sample

This repository includes a plain `database/sql` example that registers the driver with a blank import.

```bash
cd examples/driver_only
go run .
```

This reads `config.json` from the repo root and runs a simple `SELECT 1 AS ok` query against your D1 database.

### sqlx sample

The `sqlx` integration lives in its own Go module so it stays optional.

```bash
cd examples/sqlx_app
go run .
```

Or from the repo root:

```bash
go run ./examples/driver_only
```

## Status

This repository is being built incrementally using TDD and mock-based unit tests.

This project is unofficial and still experimental. It is intended for learning, prototyping, and early integration work, not production guarantees.
