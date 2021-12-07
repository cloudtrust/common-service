package database

import (
	"database/sql"

	"github.com/cloudtrust/common-service/v2/database/sqltypes"
)

type dbTransaction struct {
	tx     DbTransactionIntf
	closed bool
}

// DbTransactionIntf is exported for unit tests
type DbTransactionIntf interface {
	Commit() error
	Rollback() error

	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

// NewTransaction creates a transaction
func NewTransaction(tx DbTransactionIntf) sqltypes.Transaction {
	return &dbTransaction{tx: tx, closed: false}
}

func (tx *dbTransaction) Commit() error {
	var err = tx.tx.Commit()
	if err == nil {
		tx.closed = true
	}
	return err
}

func (tx *dbTransaction) Rollback() error {
	var err = tx.tx.Rollback()
	if err == nil {
		tx.closed = true
	}
	return err
}

func (tx *dbTransaction) Close() error {
	if tx.closed {
		return nil
	}
	return tx.Rollback()
}

func (tx *dbTransaction) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.Exec(query, args...)
}

func (tx *dbTransaction) Query(query string, args ...interface{}) (sqltypes.SQLRows, error) {
	return tx.tx.Query(query, args...)
}

func (tx *dbTransaction) QueryRow(query string, args ...interface{}) sqltypes.SQLRow {
	return tx.tx.QueryRow(query, args...)
}
