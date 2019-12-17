package sqltypes

import (
	"database/sql"
)

// SQLRows interface
type SQLRows interface {
	Next() bool
	NextResultSet() bool
	Err() error
	Scan(dest ...interface{}) error
	Close() error
}

// SQLRow interface
type SQLRow interface {
	Scan(dest ...interface{}) error
}

type sqlRowError struct {
	err error
}

// NewSQLRowError create a SQLRow based on an error
func NewSQLRowError(err error) SQLRow {
	return &sqlRowError{err: err}
}

func (sre *sqlRowError) Scan(dest ...interface{}) error {
	return sre.err
}

// CloudtrustDB interface
type CloudtrustDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (SQLRows, error)
	QueryRow(query string, args ...interface{}) SQLRow
	Ping() error
	Close() error
}

// CloudtrustDBFactory interface
type CloudtrustDBFactory interface {
	OpenDatabase() (CloudtrustDB, error)
}
