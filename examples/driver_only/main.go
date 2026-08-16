package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"

	_ "github.com/confusionhill/cf-d1-driver/pkg/d1"
)

// config contains the Cloudflare D1 credentials for a local workstation.
type config struct {
	AccountID  string `json:"account_id"`
	DatabaseID string `json:"database_id"`
	Token      string `json:"token"`
	APIBase    string `json:"api_base"`
}

func loadConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, err
	}

	var cfg config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return config{}, err
	}

	if cfg.AccountID == "" || cfg.DatabaseID == "" || cfg.Token == "" {
		return config{}, fmt.Errorf("config.json must include account_id, database_id, and token")
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

	db, err := sql.Open("d1", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Printf("ping failed: %v", err)
		return
	}

	rows, err := db.Query("SELECT 1 AS ok")
	if err != nil {
		log.Printf("query failed: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var ok int
		if err := rows.Scan(&ok); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("result: %d\n", ok)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("driver-only example completed")
}
