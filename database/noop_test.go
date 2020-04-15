package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoop(t *testing.T) {
	var cfg = GetDbConfigExt(nil, "mydb", true)

	t.Run("Database", func(t *testing.T) {
		db, _ := cfg.OpenDatabase()
		_, err := db.Query("real database would return an error")
		assert.Nil(t, err)

		db.BeginTx(nil, nil)
		db.Exec("select 1 from dual")
		db.QueryRow("select count(1) from dual").Scan()
		db.Ping()
		db.Close()
	})
	t.Run("NoopResult", func(t *testing.T) {
		var result NoopResult
		result.LastInsertId()
		result.RowsAffected()
	})
	t.Run("NoopSQLRows", func(t *testing.T) {
		var rows = NoopSQLRows{}
		defer rows.Close()

		rows.Next()
		rows.NextResultSet()
		rows.Scan()
		rows.Err()
	})
}
