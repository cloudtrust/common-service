package database

import (
	"context"
	"database/sql"

	"github.com/cloudtrust/common-service/database/sqltypes"
)

// NoopSQLRows is the result of a Query(...) for a NoopDB
type NoopSQLRows struct{}

// Next implements SQLRows.Next()
func (r *NoopSQLRows) Next() bool { return false }

// NextResultSet implements SQLRows.NextResultSet()
func (r *NoopSQLRows) NextResultSet() bool { return false }

// Err implements SQLRows.Err()
func (r *NoopSQLRows) Err() error { return nil }

// Scan implements SQLRows.Scan(...)
func (r *NoopSQLRows) Scan(dest ...interface{}) error { return nil }

// Close implements SQLRows.Close()
func (r *NoopSQLRows) Close() error { return nil }

// NoopSQLRow is the result of a QueryRow(...) for a NoopDB
type NoopSQLRow struct{}

// Scan implements SQLRow.Scan(...)
func (r *NoopSQLRow) Scan(dest ...interface{}) error {
	return nil
}

// NoopDB is a database client that does nothing.
type NoopDB struct{}

// BeginTx creates a transaction
func (db *NoopDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (sqltypes.Transaction, error) {
	return nil, nil
}

// Exec does nothing.
func (db *NoopDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return NoopResult{}, nil
}

// Query does nothing.
func (db *NoopDB) Query(query string, args ...interface{}) (sqltypes.SQLRows, error) {
	return &NoopSQLRows{}, nil
}

// QueryRow does nothing.
func (db *NoopDB) QueryRow(query string, args ...interface{}) sqltypes.SQLRow {
	return &NoopSQLRow{}
}

// Ping does nothing
func (db *NoopDB) Ping() error { return nil }

// Close does nothing
func (db *NoopDB) Close() error { return nil }

// NoopResult is a sql.Result that does nothing.
type NoopResult struct{}

// LastInsertId does nothing.
func (NoopResult) LastInsertId() (int64, error) { return 0, nil }

// RowsAffected does nothing.
func (NoopResult) RowsAffected() (int64, error) { return 0, nil }
