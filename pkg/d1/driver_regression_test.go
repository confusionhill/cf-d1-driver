package d1

import (
	"database/sql"
	"testing"
)

func TestDriverRegistered(t *testing.T) {
	db, err := sql.Open("d1", "d1://acc_123/db_456?token=secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer db.Close()
}

func TestUnsupportedTransactions(t *testing.T) {
	c := &Conn{}
	if _, err := c.Begin(); err == nil {
		t.Fatal("expected unsupported SQL transaction error")
	}
}
