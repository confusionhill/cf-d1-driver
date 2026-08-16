package d1

import "testing"

func TestBatchQueryValidate(t *testing.T) {
	q := BatchQuery{SQL: "INSERT INTO users (name) VALUES (?)", Params: []any{"alice"}}

	if err := q.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestBatchQueryValidateEmptySQL(t *testing.T) {
	q := BatchQuery{Params: []any{"alice"}}

	if err := q.Validate(); err == nil {
		t.Fatal("expected validation error for empty SQL")
	}
}
