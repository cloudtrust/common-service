package database

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/v2/database/sqltypes"

	"github.com/cloudtrust/common-service/v2/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDbVersion(t *testing.T) {
	var err error

	for _, invalidVersion := range []string{"", "1.0.1", "A.b"} {
		_, err = newDbVersion(invalidVersion)
		assert.NotNil(t, err)
	}

	var v1, v2 *dbVersion
	v1, err = newDbVersion("1.1")
	assert.Nil(t, err)

	var matchTests = map[string]bool{"0.9": false, "1.0": false, "1.1": true, "1.5": true, "2.0": true}
	for k, v := range matchTests {
		v2, _ = newDbVersion(k)
		assert.Equal(t, v, v2.matchesRequired(v1))
	}
}

func TestGetDbConnectionString(t *testing.T) {
	var conf = DbConfig{
		Username: "user",
		Password: "pass",
		Protocol: "proto",
		HostPort: "1234",
		Database: "db",
	}
	assert.Equal(t, "user:pass@proto(1234)/db", conf.getDbConnectionString())

	conf.Parameters = "params"
	assert.Equal(t, "user:pass@proto(1234)/db?params", conf.getDbConnectionString())
}

func TestConfigureDbDefault(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "mydb"
	var envUser = "the env user"
	var envPass = "the env password"

	mockConf.EXPECT().SetDefault(prefix+"-enabled", gomock.Any()).Times(1)
	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol", "-parameters", "-max-open-conns", "-max-idle-conns", "-conn-max-lifetime", "-migration", "-migration-version", "-connection-check", "-ping-timeout-ms"} {
		mockConf.EXPECT().SetDefault(prefix+suffix, gomock.Any()).Times(1)
	}
	mockConf.EXPECT().BindEnv(prefix+"-username", envUser).Times(1)
	mockConf.EXPECT().BindEnv(prefix+"-password", envPass).Times(1)

	ConfigureDbDefault(mockConf, prefix, envUser, envPass)
}

func TestGetDbConfig(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "mydb"

	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol", "-parameters"} {
		mockConf.EXPECT().GetString(prefix + suffix).Return("value" + suffix).Times(1)
	}
	for _, suffix := range []string{"-max-open-conns", "-max-idle-conns", "-conn-max-lifetime", "-ping-timeout-ms"} {
		mockConf.EXPECT().GetInt(prefix + suffix).Return(1).Times(1)
	}
	mockConf.EXPECT().GetBool(prefix + "-enabled").Return(true).Times(1)
	mockConf.EXPECT().GetBool(prefix + "-migration").Return(false).Times(1)
	mockConf.EXPECT().GetBool(prefix + "-connection-check").Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-migration-version").Return("1.0").Times(1)

	var cfg = GetDbConfig(mockConf, prefix)
	assert.Equal(t, "value-host-port", cfg.HostPort)
}

func TestCheckMigrationVersion(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var row = mock.NewSQLRow(mockCtrl)
	var dbConf = DbConfig{MigrationVersion: "1.5"}

	{
		// Can't check version: SQL query error
		var expectedError = errors.New("SQL query failed")
		var errRow = sqltypes.NewSQLRowError(expectedError)
		mockDB.EXPECT().QueryRow(gomock.Any()).Return(errRow)
		assert.Equal(t, expectedError, dbConf.checkMigrationVersion(mockDB))
	}

	{
		// Current version is higher than the minimum version requirement
		var version = "1.6"
		mockDB.EXPECT().QueryRow(gomock.Any()).Return(row)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
			var ptrVersion = dest[0].(*string)
			*ptrVersion = version
			return nil
		})
		assert.Nil(t, dbConf.checkMigrationVersion(mockDB))
	}

	{
		// Current version is higher than the minimum version requirement
		var row = mock.NewSQLRow(mockCtrl)
		var version = "1.3"
		mockDB.EXPECT().QueryRow(gomock.Any()).Return(row)
		row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
			var ptrVersion = dest[0].(*string)
			*ptrVersion = version
			return nil
		})
		var err = dbConf.checkMigrationVersion(mockDB)
		assert.NotNil(t, err)
		assert.True(t, strings.Contains(err.Error(), "not up-to-date"))
	}
}

func TestReconnectableCloudtrustDB(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var mockDBFactory = mock.NewCloudtrustDBFactory(mockCtrl)
	var expectedError = errors.New("error")

	var mockLogger = mock.NewLogger(mockCtrl)

	t.Run("Try to connect to the DB with no success", func(t *testing.T) {
		mockDBFactory.EXPECT().OpenDatabase().Return(nil, expectedError)
		_, err := NewReconnectableCloudtrustDB(mockDBFactory, mockLogger)
		assert.NotNil(t, err)
	})

	mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
	db, err := NewReconnectableCloudtrustDB(mockDBFactory, mockLogger)
	assert.Nil(t, err)

	t.Run("Exec success", func(t *testing.T) {
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, nil)
		_, err := db.Exec("request")
		assert.Nil(t, err)
	})
	t.Run("Exec failure... Ping still ok", func(t *testing.T) {
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(nil)
		_, err := db.Exec("request")
		assert.NotNil(t, err)
	})
	t.Run("Exec failure... Ping fails too...", func(t *testing.T) {
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		_, err := db.Exec("request")
		assert.NotNil(t, err)
	})

	t.Run("Query success", func(t *testing.T) {
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, nil)
		_, err := db.Query("request")
		assert.Nil(t, err)
	})
	t.Run("Query failure... Ping still ok", func(t *testing.T) {
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(nil)
		_, err := db.Query("request")
		assert.NotNil(t, err)
	})
	t.Run("Query failure... Ping fails too...", func(t *testing.T) {
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		_, err := db.Query("request")
		assert.NotNil(t, err)
	})

	t.Run("QueryRow success", func(t *testing.T) {
		var sqlRow = sqltypes.NewSQLRowError(errors.New(""))
		mockDB.EXPECT().QueryRow(gomock.Any()).Return(sqlRow)
		row := db.QueryRow("request")
		assert.Equal(t, sqlRow, row)
	})

	t.Run("Ping success", func(t *testing.T) {
		mockDB.EXPECT().Ping().Return(nil)
		assert.Nil(t, db.Ping())
	})
	t.Run("Ping failure", func(t *testing.T) {
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		assert.NotNil(t, db.Ping())
	})

	{
		db.Close()
	}

	// Wait for asynchronous reconnections
	time.Sleep(time.Second / 10)
}
