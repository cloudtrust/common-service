package database

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/viper"
)

func TestConfigureDbDefault(t *testing.T) {
	var v = viper.New()
	var prefix = "mydb"
	ConfigureDbDefault(v, prefix, "ENV_USER", "ENV_PASSWD")
	for _, suffix := range []string{"-host-port", "-username", "-password", "-database", "-protocol", "-max-open-conns", "-max-idle-conns", "-conn-max-lifetime"} {
		assert.NotNil(t, v.Get(prefix+suffix))
	}
	assert.Nil(t, v.Get("not-exits"))
}

func TestGetDbConfig(t *testing.T) {
	var v = viper.New()
	var hostport = "cloudtrust.db:3333"
	v.Set("mydb-host-port", hostport)
	var cfg = GetDbConfig(v, "mydb", false)
	assert.Equal(t, hostport, cfg.HostPort)
}

func TestOpenDatabaseNoop(t *testing.T) {
	var v = viper.New()
	var cfg = GetDbConfig(v, "mydb", true)
	db, _ := cfg.OpenDatabase()
	_, err := db.Query("real database would return an error")
	assert.Nil(t, err)
}
