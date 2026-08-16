package d1

import (
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestSQLXBindRegistration(t *testing.T) {
	// sqlx binding is a consumer concern, not a core driver dependency.
	sqlx.BindDriver("d1", sqlx.QUESTION)

	if got := sqlx.BindType("d1"); got != sqlx.QUESTION {
		t.Fatalf("expected d1 bind type to be sqlx.QUESTION, got %d", got)
	}
}

func TestSQLXOpen(t *testing.T) {
	sqlx.BindDriver("d1", sqlx.QUESTION)

	db, err := sqlx.Open("d1", "d1://acc_123/db_456?token=secret")
	if err != nil {
		t.Fatalf("unexpected sqlx open error: %v", err)
	}
	defer db.Close()
}

func TestSQLXExecUnsupportedTransaction(t *testing.T) {
	sqlx.BindDriver("d1", sqlx.QUESTION)

	db, err := sqlx.Open("d1", "d1://acc_123/db_456?token=secret")
	if err != nil {
		t.Fatalf("unexpected sqlx open error: %v", err)
	}
	defer db.Close()

	if _, err := db.Beginx(); err == nil {
		t.Fatal("expected Beginx to fail because D1 transactions are unsupported")
	}
}
