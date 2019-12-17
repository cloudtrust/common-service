package database

import "github.com/cloudtrust/common-service/database/sqltypes"

// Transaction interface
type Transaction interface {
	Commit() error
	Rollback() error

	// Close: if not explicitely Commited or Rolled back, Rollback the transaction
	Close() error
}

type dbTransaction struct {
	db     sqltypes.CloudtrustDB
	closed bool
}

// NewTransaction creates a transaction
func NewTransaction(db sqltypes.CloudtrustDB) (Transaction, error) {
	var _, err = db.Exec("START TRANSACTION")
	if err != nil {
		return nil, err
	}
	return &dbTransaction{db: db, closed: false}, nil
}

func (tx *dbTransaction) close(cmd string) error {
	if tx.closed {
		return nil
	}
	var _, err = tx.db.Exec(cmd)
	tx.closed = true
	return err
}

func (tx *dbTransaction) Commit() error {
	return tx.close("COMMIT")
}

func (tx *dbTransaction) Rollback() error {
	return tx.close("ROLLBACK")
}

func (tx *dbTransaction) Close() error {
	return tx.Rollback()
}
