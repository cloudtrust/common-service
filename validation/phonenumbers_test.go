package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObfuscatePhoneNumber(t *testing.T) {
	t.Run("Invalid phone number-too small to obfuscate", func(t *testing.T) {
		assert.Equal(t, "+12", ObfuscatePhoneNumber("+12"))
	})
	t.Run("Invalid phone number-keep 2 digits on both sides of phone number", func(t *testing.T) {
		assert.Equal(t, "ZY1**77", ObfuscatePhoneNumber("ZY18877"))
	})
	t.Run("Valid phone number-Keep country code and last 2 digits", func(t *testing.T) {
		assert.Equal(t, "+33*******99", ObfuscatePhoneNumber("+33686115599"))
	})
}
