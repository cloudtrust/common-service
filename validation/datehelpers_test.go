package validation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLargeDuration(t *testing.T) {
	var now = time.Now()
	t.Run("Empty duration", func(t *testing.T) {
		assert.False(t, IsValidLargeDuration(""))
		assert.Equal(t, now, AddLargeDuration(now, ""))
	})
	t.Run("Invalid duration", func(t *testing.T) {
		var duration = "2d3y4"
		assert.False(t, IsValidLargeDuration(duration))
	})
	t.Run("Valid duration", func(t *testing.T) {
		var duration = "2d3y4w1m"
		assert.True(t, IsValidLargeDuration(duration))

		var after = now.AddDate(3, 1, 4*7+2)
		assert.Equal(t, after, AddLargeDuration(now, duration))
	})
}
