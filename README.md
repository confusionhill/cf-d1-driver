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

For `sqlx`, bind the driver and import the package normally:

```go
import (
    "github.com/jmoiron/sqlx"
    d1 "github.com/confusionhill/cf-d1-driver/pkg/d1"
)

func init() {
    sqlx.BindDriver("d1", sqlx.QUESTION)
}
```

## Transaction model

D1 does not support SQL `BEGIN`, `COMMIT`, or `ROLLBACK` statements.

Use the batch API for atomic multi-step writes:

```go
batch := []d1.BatchQuery{
    {SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}},
    {SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"alicia", "alice"}},
}
```

## Sample app

A sample executable lives under `examples/sqlx_app` and demonstrates `sqlx` usage with the D1 driver.

```bash
go run ./examples/sqlx_app
```

## Status

This repository is being built incrementally using TDD and mock-based unit tests.

This project is unofficial and still experimental. It is intended for learning, prototyping, and early integration work, not production guarantees.

## Transaction model

D1 does not support SQL `BEGIN`, `COMMIT`, or `ROLLBACK` statements.

Use the batch API for atomic multi-step writes:

```go
batch := []d1.BatchQuery{
    {SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}},
    {SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"alicia", "alice"}},
}
```

## Sample app

A sample executable lives under `examples/sqlx_app` and demonstrates `sqlx` usage with the D1 driver.

```bash
go run ./examples/sqlx_app
```

## Status

This repository is being built incrementally using TDD and mock-based unit tests.

This project is unofficial and still experimental. It is intended for learning, prototyping, and early integration work, not production guarantees.
