package http

import (
	"errors"
	"testing"

	"github.com/rs/cors"
	"github.com/stretchr/testify/assert"
)

func assertEmptyCors(t *testing.T, cors cors.Options) {
	assert.Nil(t, cors.AllowedOrigins)
	assert.Nil(t, cors.AllowedMethods)
	assert.Equal(t, false, cors.AllowCredentials)
	assert.Nil(t, cors.AllowedHeaders)
	assert.Nil(t, cors.ExposedHeaders)
	assert.Equal(t, false, cors.Debug)
}

func unmarshalError(key string, rawVal any) error {
	return errors.New("Failed")
}

func unmarshalSuccess(key string, rawVal any) error {
	var ptr *CorsConfiguration = rawVal.(*CorsConfiguration)
	ptr.AllowedOrigins = []string{"origin1", "origin2"}
	ptr.AllowedMethods = []string{"method1", "method2"}
	ptr.AllowCredentials = true
	ptr.AllowedHeaders = []string{"header11", "header12"}
	ptr.ExposedHeaders = []string{"header21", "header22"}
	ptr.Debug = true
	return nil
}

func TestGetCorsOptions(t *testing.T) {
	t.Run("Unmarshalling fails", func(t *testing.T) {
		cors, err := GetCorsOptions("myconf", unmarshalError)
		assert.NotNil(t, err)
		assertEmptyCors(t, cors)
	})
	t.Run("Unmarshal success", func(t *testing.T) {
		cors, err := GetCorsOptions("myconf", unmarshalSuccess)
		assert.Nil(t, err)
		assert.Len(t, cors.AllowedOrigins, 2)
		assert.Len(t, cors.AllowedMethods, 2)
		assert.True(t, cors.AllowCredentials)
		assert.Len(t, cors.AllowedHeaders, 2)
		assert.Len(t, cors.ExposedHeaders, 2)
		assert.True(t, cors.Debug)
	})
}

type configMock struct {
}

func (c *configMock) GetBool(key string) bool {
	return key == "myconf-allow-credential" || key == "myconf-debug"
}

func (c *configMock) GetStringSlice(key string) []string {
	if key == "myconf-allowed-origins" || key == "myconf-allowed-methods" || key == "myconf-allowed-headers" || key == "myconf-exposed-headers" {
		return []string{"value1", "value2"}
	}
	return nil
}

func TestGetCorsOptionsLegacy(t *testing.T) {
	var config = &configMock{}

	t.Run("Key not configured", func(t *testing.T) {
		cors, err := GetCorsOptionsLegacy(config, "not-configured", unmarshalError)
		assert.NotNil(t, err)
		assertEmptyCors(t, cors)
	})
	t.Run("Configured with struct", func(t *testing.T) {
		cors, err := GetCorsOptionsLegacy(config, "myConf", unmarshalSuccess)
		assert.Nil(t, err)
		assert.Len(t, cors.AllowedOrigins, 2)
		assert.Len(t, cors.AllowedMethods, 2)
		assert.True(t, cors.AllowCredentials)
		assert.Len(t, cors.AllowedHeaders, 2)
		assert.Len(t, cors.ExposedHeaders, 2)
		assert.True(t, cors.Debug)
	})
	t.Run("Use legacy configuration", func(t *testing.T) {
		cors, err := GetCorsOptionsLegacy(config, "myconf", unmarshalError)
		assert.Nil(t, err)
		assert.Len(t, cors.AllowedOrigins, 2)
		assert.Len(t, cors.AllowedMethods, 2)
		assert.True(t, cors.AllowCredentials)
		assert.Len(t, cors.AllowedHeaders, 2)
		assert.Len(t, cors.ExposedHeaders, 2)
		assert.True(t, cors.Debug)
	})
}
