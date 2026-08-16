package d1

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestSQLXOpen(t *testing.T) {
	db, err := sqlx.Open("d1", "d1://acc_123/db_456?token=secret")
	if err != nil {
		t.Fatalf("unexpected sqlx open error: %v", err)
	}
	defer db.Close()
}

func TestSQLXExecUnsupportedTransaction(t *testing.T) {
	db, err := sqlx.Open("d1", "d1://acc_123/db_456?token=secret")
	if err != nil {
		t.Fatalf("unexpected sqlx open error: %v", err)
	}
	defer db.Close()

	if _, err := db.Beginx(); err == nil {
		t.Fatal("expected Beginx to fail because D1 transactions are unsupported")
	}
}
