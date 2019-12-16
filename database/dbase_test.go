package database

import (
	"testing"
	"time"

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

func TestOpenDatabaseNoop(t *testing.T) {
	var cfg = GetDbConfig(nil, "mydb", true)
	db, _ := cfg.OpenDatabase()
	_, err := db.Query("real database would return an error")
	assert.Nil(t, err)

	db.Exec("select 1 from dual")
	db.QueryRow("select count(1) from dual")
	db.Ping()
	db.SetConnMaxLifetime(time.Duration(1))
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}

func TestNoopResult(t *testing.T) {
	var result NoopResult
	result.LastInsertId()
	result.RowsAffected()
}
