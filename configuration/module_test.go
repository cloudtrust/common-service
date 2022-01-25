package configuration

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/cloudtrust/common-service/v2/configuration/mock"
	"github.com/cloudtrust/common-service/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetRealmConfigurations(t *testing.T) {
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
		mockSQLRow.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(errors.New("SQL error"))
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("SQL No row", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any(), gomock.Any()).Return(sql.ErrNoRows)
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Invalid JSON", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
			*(dest[0]).(*string) = `{`
			*(dest[1]).(*string) = `{}`
			return nil
		})
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.NotNil(t, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockDB.EXPECT().QueryRow(gomock.Any(), realmID).Return(mockSQLRow)
		mockSQLRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
			*(dest[0]).(*string) = `{}`
			*(dest[1]).(*string) = `{}`
			return nil
		})
		var _, _, err = module.GetRealmConfigurations(ctx, realmID)
		assert.Nil(t, err)
	})
}

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

func TestGetAuthorizations(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var mockSQLRows = mock.NewSQLRows(mockCtrl)
	var logger = log.NewNopLogger()

	var notInScopeAction = "ACTION#1"
	var allowedAction = "ACTION#2"
	var actions = []string{allowedAction}
	var ctx = context.TODO()

	var module = NewConfigurationReaderDBModule(mockDB, logger, actions)

	t.Run("Query fails", func(t *testing.T) {
		var sqlError = errors.New("SQL error")
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, sqlError)

		var _, err = module.GetAuthorizations(ctx)
		assert.Equal(t, sqlError, err)
	})

	// Now, query will always be successful
	mockDB.EXPECT().Query(gomock.Any()).Return(mockSQLRows, nil).AnyTimes()
	mockSQLRows.EXPECT().Close().AnyTimes()

	t.Run("scan fails", func(t *testing.T) {
		var scanError = errors.New("scan error")
		mockSQLRows.EXPECT().Next().Return(true)
		mockSQLRows.EXPECT().Scan(gomock.Any()).Return(scanError)

		var _, err = module.GetAuthorizations(ctx)
		assert.Equal(t, scanError, err)
	})

	t.Run("Query returns 2 rows", func(t *testing.T) {
		gomock.InOrder(
			mockSQLRows.EXPECT().Next().Return(true),
			mockSQLRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
				*(dest[0]).(*string) = "realm#1"
				*(dest[1]).(*string) = "group#1"
				*(dest[2]).(*string) = notInScopeAction
				*(dest[3]).(*sql.NullString) = sql.NullString{Valid: false}
				*(dest[4]).(*sql.NullString) = sql.NullString{Valid: false}
				return nil
			}),
			mockSQLRows.EXPECT().Next().Return(true),
			mockSQLRows.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
				*(dest[0]).(*string) = "realm#2"
				*(dest[1]).(*string) = "group#2"
				*(dest[2]).(*string) = allowedAction
				*(dest[3]).(*sql.NullString) = sql.NullString{Valid: true, String: "targetRealm"}
				*(dest[4]).(*sql.NullString) = sql.NullString{Valid: true, String: "targetGroup"}
				return nil
			}),
			mockSQLRows.EXPECT().Next().Return(false),
		)

		var res, err = module.GetAuthorizations(ctx)
		assert.Nil(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, allowedAction, *res[0].Action)
	})
}

func TestIsInScope(t *testing.T) {
	var module = NewConfigurationReaderDBModule(nil, nil)
	assert.True(t, module.isInAuthorizationScope("any-auth-will-be-considered-in-scope"))

	module = NewConfigurationReaderDBModule(nil, nil, []string{"autun"}, []string{"auth2", "auth3"})
	assert.False(t, module.isInAuthorizationScope("any-auth-will-be-considered-in-scope"))
	assert.True(t, module.isInAuthorizationScope("autun"))
	assert.True(t, module.isInAuthorizationScope("auth2"))
	assert.True(t, module.isInAuthorizationScope("auth3"))
}
