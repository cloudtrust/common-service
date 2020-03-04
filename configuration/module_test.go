package configuration

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/configuration/mock"
	"github.com/cloudtrust/common-service/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var mockSQLRow = mock.NewSQLRow(mockCtrl)
	var logger = log.NewNopLogger()

	var realmID = "myrealm"
	var ctx = context.TODO()
	var module = NewConfigurationReaderDBModule(mockDB, logger)

	t.Run("SQL error", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).Return(errors.New("SQL error"))
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(ptrConfig *string) error {
			*ptrConfig = `{}`
			return nil
		})
		var _, err = module.GetConfiguration(ctx, realmID)
		assert.Nil(t, err)
	})
}

func TestGetAdminConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var mockSQLRow = mock.NewSQLRow(mockCtrl)
	var logger = log.NewNopLogger()

	var realmID = "myrealm"
	var ctx = context.TODO()
	var module = NewConfigurationReaderDBModule(mockDB, logger)

	t.Run("SQL error", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).Return(errors.New("SQL error"))
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).Return(sql.ErrNoRows)
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(ptrConfig *string) error {
			*ptrConfig = `{}`
			return nil
		})
		var _, err = module.GetAdminConfiguration(ctx, realmID)
		assert.Nil(t, err)
	})
}
