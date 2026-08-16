package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"

	d1 "github.com/dika/cf-d1-go-driver/pkg/d1"
)

func init() {
	sqlx.BindDriver("d1", sqlx.QUESTION)
}

type appConfig struct {
	AccountID  string `json:"account_id"`
	DatabaseID string `json:"database_id"`
	Token      string `json:"token"`
	APIBase    string `json:"api_base"`
}

func loadConfig(path string) (appConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return appConfig{}, err
	}

	var cfg appConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return appConfig{}, err
	}

	if cfg.AccountID == "" || cfg.DatabaseID == "" || cfg.Token == "" {
		return appConfig{}, fmt.Errorf("config.json must include account_id, database_id, and token")
	}
	if cfg.APIBase == "" {
		cfg.APIBase = "https://api.cloudflare.com/client/v4"
	}

	return cfg, nil
}

func main() {
	cfg, err := loadConfig("config.json")
	if err != nil {
		log.Fatal("load config.json: ", err)
	}

	dsn := fmt.Sprintf(
		"d1://%s/%s?token=%s&api_base=%s",
		url.PathEscape(cfg.AccountID),
		url.PathEscape(cfg.DatabaseID),
		url.QueryEscape(cfg.Token),
		url.QueryEscape(cfg.APIBase),
	)

	db, err := sqlx.Open("d1", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// This is just a sample app. The actual D1 HTTP layer is still scaffolded,
	// but the design matches Cloudflare D1 semantics and sqlx usage patterns.
	_, err = db.Exec("SELECT 1")
	if err != nil {
		fmt.Println("driver is scaffolded, execution depends on valid D1 credentials:", err)
	}

	batch := []d1.BatchQuery{
		{SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}},
		{SQL: "UPDATE users SET name = ? WHERE name = ?", Params: []any{"alicia", "alice"}},
	}

	_ = batch
	fmt.Println("sample app initialized")
}
