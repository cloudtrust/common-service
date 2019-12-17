package database

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/cloudtrust/common-service/database/sqltypes"

	"github.com/cloudtrust/common-service/database/mock"
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
		v2, err = newDbVersion(k)
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

	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol", "-parameters", "-max-open-conns", "-max-idle-conns", "-conn-max-lifetime", "-migration", "-migration-version", "-connection-check"} {
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
	for _, suffix := range []string{"-max-open-conns", "-max-idle-conns", "-conn-max-lifetime"} {
		mockConf.EXPECT().GetInt(prefix + suffix).Return(1).Times(1)
	}
	mockConf.EXPECT().GetBool(prefix + "-migration").Return(false).Times(1)
	mockConf.EXPECT().GetBool(prefix + "-connection-check").Return(true).Times(1)
	mockConf.EXPECT().GetString(prefix + "-migration-version").Return("1.0").Times(1)

	var cfg = GetDbConfig(mockConf, prefix, false)
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

func TestOpenDatabaseNoop(t *testing.T) {
	var cfg = GetDbConfig(nil, "mydb", true)
	db, _ := cfg.OpenDatabase()
	_, err := db.Query("real database would return an error")
	assert.Nil(t, err)

	db.Exec("select 1 from dual")
	db.QueryRow("select count(1) from dual")
	db.Ping()
	db.Close()
}

func TestNoopResult(t *testing.T) {
	var result NoopResult
	result.LastInsertId()
	result.RowsAffected()
}

func TestReconnectableCloudtrustDB(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockDB = mock.NewCloudtrustDB(mockCtrl)
	var mockDBFactory = mock.NewCloudtrustDBFactory(mockCtrl)
	var expectedError = errors.New("error")

	{
		// Try to connect to the DB with no success
		mockDBFactory.EXPECT().OpenDatabase().Return(nil, expectedError)
		_, err := NewReconnectableCloudtrustDB(mockDBFactory)
		assert.NotNil(t, err)
	}

	// Get a connection to the DB
	mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
	db, err := NewReconnectableCloudtrustDB(mockDBFactory)
	assert.Nil(t, err)

	{
		// Exec success
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, nil)
		_, err := db.Exec("request")
		assert.Nil(t, err)

		// Exec failure... Ping still ok
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(nil)
		_, err = db.Exec("request")
		assert.NotNil(t, err)

		// Exec failure... Ping fails too...
		mockDB.EXPECT().Exec(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		_, err = db.Exec("request")
		assert.NotNil(t, err)
	}

	{
		// Query success
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, nil)
		_, err := db.Query("request")
		assert.Nil(t, err)

		// Query failure... Ping still ok
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(nil)
		_, err = db.Query("request")
		assert.NotNil(t, err)

		// Query failure... Ping fails too...
		mockDB.EXPECT().Query(gomock.Any()).Return(nil, expectedError)
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		_, err = db.Query("request")
		assert.NotNil(t, err)
	}

	{
		var sqlRow = sqltypes.NewSQLRowError(errors.New(""))

		// QueryRow success
		mockDB.EXPECT().QueryRow(gomock.Any()).Return(sqlRow)
		row := db.QueryRow("request")
		assert.Equal(t, sqlRow, row)
	}

	{
		// Ping success
		mockDB.EXPECT().Ping().Return(nil)
		assert.Nil(t, db.Ping())

		// Ping failure
		mockDB.EXPECT().Ping().Return(expectedError)
		mockDB.EXPECT().Close()
		mockDBFactory.EXPECT().OpenDatabase().Return(mockDB, nil)
		assert.NotNil(t, db.Ping())
	}

	{
		db.Close()
	}

	// Wait for asynchronous reconnections
	time.Sleep(time.Second / 10)
}
