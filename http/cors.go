package http

import (
	"errors"
	"strings"

	cs "github.com/cloudtrust/common-service/v2"
	"github.com/rs/cors"
)

// Configuration interface
type Configuration interface {
	GetBool(key string) bool
	GetStringSlice(key string) []string
}

// CorsConfiguration struct
type CorsConfiguration struct {
	AllowedOrigins   []string `mapstructure:"allowed-origins"`
	AllowedMethods   []string `mapstructure:"allowed-methods"`
	AllowCredentials bool     `mapstructure:"allow-credential"`
	AllowedHeaders   []string `mapstructure:"allowed-headers"`
	ExposedHeaders   []string `mapstructure:"exposed-headers"`
	Debug            bool     `mapstructure:"debug"`
}

func (c *CorsConfiguration) isEmpty() bool {
	return len(c.AllowedOrigins) == 0 && len(c.AllowedMethods) == 0 && !c.AllowCredentials && len(c.AllowedHeaders) == 0 && len(c.ExposedHeaders) == 0 && !c.Debug

}

func (c *CorsConfiguration) toCorsOptions() cors.Options {
	return cors.Options{
		AllowedOrigins:   c.AllowedOrigins,
		AllowedMethods:   c.AllowedMethods,
		AllowCredentials: c.AllowCredentials,
		AllowedHeaders:   c.AllowedHeaders,
		ExposedHeaders:   c.ExposedHeaders,
		Debug:            c.Debug,
	}
}

// ConfigureCorsDefault configures default CORS values
func ConfigureCorsDefault(v cs.Configuration, name string) {
	v.SetDefault(name+".allowed-origins", []string{})
	v.SetDefault(name+".allowed-methods", []string{})
	v.SetDefault(name+".allow-credentials", true)
	v.SetDefault(name+".allowed-headers", []string{})
	v.SetDefault(name+".exposed-headers", []string{})
	v.SetDefault(name+".debug", false)
}

// GetCorsOptions returns CORS options from configuration
func GetCorsOptions(name string, unmarshallFunc func(key string, rawVal any) error) (cors.Options, error) {
	var corsConf CorsConfiguration
	if err := unmarshallFunc(name, &corsConf); err != nil {
		return cors.Options{}, err
	}
	return cors.Options{
		AllowedOrigins:   corsConf.AllowedOrigins,
		AllowedMethods:   corsConf.AllowedMethods,
		AllowCredentials: corsConf.AllowCredentials,
		AllowedHeaders:   corsConf.AllowedHeaders,
		ExposedHeaders:   corsConf.ExposedHeaders,
		Debug:            corsConf.Debug,
	}, nil
}

// GetCorsOptionsLegacy returns CORS options from configuration
func GetCorsOptionsLegacy(c Configuration, name string, unmarshalFunc func(key string, rawVal any) error) (cors.Options, error) {
	var providedName = name
	if corsConf, err := GetCorsOptions(name, unmarshalFunc); err == nil {
		return corsConf, nil
	}
	if !strings.HasSuffix(name, "-") {
		name += "-"
	}
	var res = CorsConfiguration{
		AllowedOrigins:   c.GetStringSlice(name + "allowed-origins"),
		AllowedMethods:   c.GetStringSlice(name + "allowed-methods"),
		AllowCredentials: c.GetBool(name + "allow-credential"),
		AllowedHeaders:   c.GetStringSlice(name + "allowed-headers"),
		ExposedHeaders:   c.GetStringSlice(name + "exposed-headers"),
		Debug:            c.GetBool(name + "debug"),
	}
	if res.isEmpty() {
		return cors.Options{}, errors.New("CORS not configured with key " + providedName)
	}
	return res.toCorsOptions(), nil
}
