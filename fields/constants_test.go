package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKnownFields(t *testing.T) {
	assert.NotEmpty(t, GetKnownFields())
}

func TestField(t *testing.T) {
	var (
		key   = "the key"
		attrb = "the attribute"

		f Field = &field{
			Field:     key,
			Attribute: attrb,
		}
	)
	assert.Equal(t, key, f.Key())
	assert.Equal(t, attrb, f.AttributeName())
}
