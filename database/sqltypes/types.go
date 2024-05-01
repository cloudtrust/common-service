package sqltypes

import (
	"context"
	"database/sql"
)

// SQLRows interface
type SQLRows interface {
	Next() bool
	NextResultSet() bool
	Err() error
	Scan(dest ...any) error
	Close() error
}

// SQLRow interface
type SQLRow interface {
	Scan(dest ...any) error
}

type sqlRowError struct {
	err error
}

// NewSQLRowError create a SQLRow based on an error
func NewSQLRowError(err error) SQLRow {
	return &sqlRowError{err: err}
}

func (sre *sqlRowError) Scan(dest ...any) error {
	return sre.err
}

// CloudtrustDB interface
type CloudtrustDB interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Transaction, error)
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (SQLRows, error)
	QueryRow(query string, args ...any) SQLRow
	Ping() error
	Close() error
	Stats() sql.DBStats
}

// CloudtrustDBFactory interface
type CloudtrustDBFactory interface {
	OpenDatabase() (CloudtrustDB, error)
}

// Transaction interface
type Transaction interface {
	Commit() error
	Rollback() error

	// Close: if not explicitely Commited or Rolled back, Rollback the transaction
	Close() error

	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (SQLRows, error)
	QueryRow(query string, args ...any) SQLRow
}
