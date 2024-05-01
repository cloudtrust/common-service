package commonservice

import (
	"context"
	"time"
)

// Endpoint type
type Endpoint func(ctx context.Context, request any) (response any, err error)

// Middleware type
type Middleware func(Endpoint) Endpoint

// Configuration interface
type Configuration interface {
	SetDefault(key string, value any)

	BindEnv(input ...string) error

	Set(key string, value any)

	Get(key string) any
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
