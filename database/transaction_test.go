package database

import (
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTransaction(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockTx = mock.NewDbTransactionIntf(mockCtrl)
	var sqlError = errors.New("I'm a SQL error")
	var query = "select columns from table"
	var param1 = "param1"

	t.Run("Exec", func(t *testing.T) {
		var tx = NewTransaction(mockTx)
		defer tx.Close()

		mockTx.EXPECT().Exec(query, param1).Return(nil, sqlError)
		mockTx.EXPECT().Rollback().Return(nil)
		var _, err = tx.Exec(query, param1)
		assert.Equal(t, sqlError, err)
	})
	t.Run("Query", func(t *testing.T) {
		var tx = NewTransaction(mockTx)
		defer tx.Close()

		mockTx.EXPECT().Query(query, param1).Return(nil, sqlError)
		mockTx.EXPECT().Rollback().Return(nil)
		var _, err = tx.Query(query, param1)
		assert.Equal(t, sqlError, err)
		// Force rollback... tx.Close() won't have to do it
		tx.Rollback()
	})
	t.Run("QueryRow", func(t *testing.T) {
		var tx = NewTransaction(mockTx)
		defer tx.Close()

		mockTx.EXPECT().QueryRow(query, param1).Return(nil)
		mockTx.EXPECT().Commit().Return(nil)
		assert.Nil(t, tx.QueryRow(query, param1))
		// Force commit... tx.Close() won't have to rollback
		tx.Commit()
	})
}
