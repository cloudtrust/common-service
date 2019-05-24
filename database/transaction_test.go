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

	var tx Transaction
	var mockDB = mock.NewCloudtrustDB(mockCtrl)

	mockDB.EXPECT().Exec("START TRANSACTION").Return(nil, errors.New("db error")).Times(1)
	_, err := NewTransaction(mockDB)
	assert.NotNil(t, err)

	mockDB.EXPECT().Exec("START TRANSACTION").Return(nil, nil).Times(1)
	tx, err = NewTransaction(mockDB)
	assert.Nil(t, err)

	mockDB.EXPECT().Exec("COMMIT").Return(nil, nil).Times(1)
	assert.Nil(t, tx.Commit())

	// Already closed
	assert.Nil(t, tx.Close())
}
