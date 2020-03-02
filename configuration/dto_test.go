package configuration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRealmConfiguration(t *testing.T) {
	t.Run("Invalid JSON", func(t *testing.T) {
		var _, err = NewRealmConfiguration(`{`)
		assert.NotNil(t, err)
	})
	t.Run("Has deprecated field, new field not set", func(t *testing.T) {
		var conf, _ = NewRealmConfiguration(`{"api_self_mail_editing_enabled":true}`)
		assert.Nil(t, conf.DeprecatedAPISelfMailEditingEnabled)
		assert.True(t, *conf.APISelfAccountEditingEnabled)
	})
	t.Run("Has deprecated field, new field already set", func(t *testing.T) {
		var conf, _ = NewRealmConfiguration(`{"api_self_mail_editing_enabled":true, "api_self_account_editing_enabled":false}`)
		assert.Nil(t, conf.DeprecatedAPISelfMailEditingEnabled)
		assert.False(t, *conf.APISelfAccountEditingEnabled)
	})
}

func TestNewRealmAdminConfiguration(t *testing.T) {
	t.Run("Invalid JSON", func(t *testing.T) {
		var _, err = NewRealmAdminConfiguration(`{`)
		assert.NotNil(t, err)
	})
	t.Run("Valid JSON", func(t *testing.T) {
		var _, err = NewRealmAdminConfiguration(`{}`)
		assert.Nil(t, err)
	})
}
