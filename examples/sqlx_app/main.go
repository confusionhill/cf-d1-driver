package main

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"

    d1 "github.com/dika/cf-d1-go-driver/pkg/d1"
)

func main() {
    db, err := sqlx.Open("d1", "d1://acc_123/db_456?token=secret")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // This is just a sample app. The actual D1 HTTP layer is still scaffolded,
    // but the design matches cloudflare D1 semantics and sqlx usage patterns.
    _, err = db.Exec("SELECT 1")
    if err != nil {
        fmt.Println("driver is scaffolded, execution depends on real D1 credentials:", err)
    }

    batch := []d1.BatchQuery{
        {SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}},
        {SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"alicia", "alice"}},
    }

    _ = batch
    fmt.Println("sample app initialized")
}
