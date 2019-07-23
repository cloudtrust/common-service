package commonservice

import (
	"context"
	"time"
)

// Endpoint type
type Endpoint func(ctx context.Context, request interface{}) (response interface{}, err error)

// Middleware type
type Middleware func(Endpoint) Endpoint

// Configuration interface
type Configuration interface {
	SetDefault(key string, value interface{})

	BindEnv(input ...string) error

	Set(key string, value interface{})

	Get(key string) interface{}
	GetString(key string) string
	GetStringSlice(key string) []string
	GetBool(key string) bool
	GetInt(key string) int
	GetInt32(key string) int32
	GetInt64(key string) int64
	GetFloat64(key string) float64
	GetTime(key string) time.Time
	GetDuration(key string) time.Duration
}
