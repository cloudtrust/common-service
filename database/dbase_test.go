package database

import (
	"testing"
	"time"

	"github.com/cloudtrust/common-service/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDbVersion(t *testing.T) {
	_, err := newDbVersion("")
	assert.NotNil(t, err)

	_, err = newDbVersion("1.0.1")
	assert.NotNil(t, err)

	_, err = newDbVersion("A.b")
	assert.NotNil(t, err)

	var v1, v2 *dbVersion
	v1, err = newDbVersion("1.0")
	assert.Nil(t, err)

	v2, err = newDbVersion("0.9")
	assert.Nil(t, err)
	assert.False(t, v2.matchesRequired(v1))

	v2, err = newDbVersion("1.0")
	assert.Nil(t, err)
	assert.True(t, v2.matchesRequired(v1))

	v2, err = newDbVersion("1.5")
	assert.Nil(t, err)
	assert.True(t, v2.matchesRequired(v1))

	v2, err = newDbVersion("2.0")
	assert.Nil(t, err)
	assert.True(t, v2.matchesRequired(v1))
}

func TestConfigureDbDefault(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockConf = mock.NewConfiguration(mockCtrl)

	var prefix = "mydb"
	var envUser = "the env user"
	var envPass = "the env password"

	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol", "-max-open-conns", "-max-idle-conns", "-conn-max-lifetime"} {
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

	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol"} {
		mockConf.EXPECT().GetString(prefix + suffix).Return("value" + suffix).Times(1)
	}
	for _, suffix := range []string{"-max-open-conns", "-max-idle-conns", "-conn-max-lifetime"} {
		mockConf.EXPECT().GetInt(prefix + suffix).Return(1).Times(1)
	}

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
	db.SetConnMaxLifetime(time.Duration(1))
	db.SetMaxIdleConns(1)
	db.SetMaxOpenConns(1)
}

func TestNoopResult(t *testing.T) {
	var result NoopResult
	result.LastInsertId()
	result.RowsAffected()
}
