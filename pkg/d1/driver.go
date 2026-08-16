package d1

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"sync"

	cf "github.com/dika/cf-d1-go-driver/internal/cloudflare"
)

func init() {
	sql.Register("d1", &Driver{})
}

// DB is the abstraction boundary for the driver package.
// This keeps the implementation testable and allows mock-backed adapters.
type DB interface {
	ExecContext(ctx context.Context, query string, args ...any) (driver.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (driver.Rows, error)
	Batch(ctx context.Context, queries ...BatchQuery) ([]BatchResult, error)
	Ping(ctx context.Context) error
}

// QueryExecutor describes the subset of behavior needed for sqlx compatibility.
type QueryExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (driver.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (driver.Rows, error)
}

// Driver is the exported driver type used by database/sql.
type Driver struct{}

// Open creates a new sql.DB for the D1 driver.
func Open(dsn string) (*sql.DB, error) {
	return sql.Open("d1", dsn)
}

// Conn is the concrete connection used by database/sql.
type Conn struct {
	cfg    Config
	client cf.Client
	mu     sync.Mutex
}

// OpenConnector is used for custom connector integration.
func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	cfg, err := ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	return &Connector{cfg: cfg, client: cf.NewClient(cfg.AccountID, cfg.DatabaseID, cfg.APIToken, cfg.APIBaseURL)}, nil
}

// Open implements driver.Driver.
func (d *Driver) Open(name string) (driver.Conn, error) {
	cfg, err := ParseDSN(name)
	if err != nil {
		return nil, err
	}
	return &Conn{cfg: cfg, client: cf.NewClient(cfg.AccountID, cfg.DatabaseID, cfg.APIToken, cfg.APIBaseURL)}, nil
}

// Connector provides connector-based connection creation.
type Connector struct {
	cfg    Config
	client cf.Client
}

func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	return &Conn{cfg: c.cfg, client: c.client}, nil
}

func (c *Connector) Driver() driver.Driver {
	return &Driver{}
}

// Prepare is required by database/sql/driver.Conn.
func (c *Conn) Prepare(query string) (driver.Stmt, error) {
	return nil, fmt.Errorf("d1: prepared statements are not supported; use ExecContext / QueryContext")
}

// Close closes the connection.
func (c *Conn) Close() error { return nil }

// Begin is intentionally unsupported because D1 does not support SQL transactions.
func (c *Conn) Begin() (driver.Tx, error) {
	return nil, fmt.Errorf("d1: SQL transactions are unsupported; use Batch() for atomic multi-statement operations")
}

// ExecContext executes a single D1 query.
func (c *Conn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	params := make([]any, len(args))
	for i, a := range args {
		params[i] = a.Value
	}

	results, err := c.client.Query(ctx, cf.QueryRequest{SQL: query, Params: params})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return driver.RowsAffected(0), nil
	}
	return &Result{rowsAffected: results[0].Meta.Changes, lastInsertID: results[0].Meta.LastRowID}, nil
}

// QueryContext executes a D1 query and returns rows.
func (c *Conn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	params := make([]any, len(args))
	for i, a := range args {
		params[i] = a.Value
	}

	results, err := c.client.Query(ctx, cf.QueryRequest{SQL: query, Params: params})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return &Rows{columns: []string{}, rows: nil}, nil
	}

	rows := make([]map[string]any, 0, len(results[0].Results))
	for _, item := range results[0].Results {
		mr := make(map[string]any, len(item))
		for k, v := range item {
			mr[k] = v
		}
		rows = append(rows, mr)
	}

	columns := make([]string, 0)
	for _, row := range rows {
		for k := range row {
			columns = append(columns, k)
		}
		if len(columns) > 0 {
			break
		}
	}

	return &Rows{columns: columns, rows: rows, index: -1}, nil
}

// Ping checks the D1 endpoint minimally.
func (c *Conn) Ping(ctx context.Context) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("d1: client is not initialized")
	}
	_, err := c.client.Query(ctx, cf.QueryRequest{SQL: "SELECT 1"})
	return err
}

// Result implements driver.Result.
type Result struct {
	rowsAffected int64
	lastInsertID int64
}

func (r *Result) LastInsertId() (int64, error) { return r.lastInsertID, nil }
func (r *Result) RowsAffected() (int64, error) { return r.rowsAffected, nil }

// Rows implements driver.Rows.
type Rows struct {
	columns []string
	rows    []map[string]any
	index   int
}

func (r *Rows) Columns() []string { return r.columns }
func (r *Rows) Close() error { return nil }
func (r *Rows) Next(dest []driver.Value) error {
	r.index++
	if r.index >= len(r.rows) {
		return io.EOF
	}
	row := r.rows[r.index]
	for i, col := range r.columns {
		if val, ok := row[col]; ok {
			dest[i] = val
		} else {
			dest[i] = nil
		}
	}
	return nil
}

// Check interfaces compile.
var _ driver.Driver = (*Driver)(nil)
var _ driver.Connector = (*Connector)(nil)
var _ driver.Conn = (*Conn)(nil)
var _ driver.Result = (*Result)(nil)
var _ driver.Rows = (*Rows)(nil)
