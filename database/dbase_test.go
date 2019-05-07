package database

//go:generate mockgen -destination=./mock/configuration.go -package=mock -mock_names=Configuration=Configuration github.com/cloudtrust/common-service Configuration

import (
	"testing"
	"time"

	"github.com/cloudtrust/common-service/database/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

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
