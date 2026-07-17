package middleware

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBasicAuthFromHeader(t *testing.T) {
	bac := NewBasicAuthCollection()

	t.Run("Valid basic auth", func(t *testing.T) {
		basic, err := bac.getBasicAuthFromHeader("Basic     dXNlcjpwYXNzMA==   ")
		assert.Nil(t, err)
		assert.Equal(t, "dXNlcjpwYXNzMA==", basic)
	})

	t.Run("Not a basic auth", func(t *testing.T) {
		_, err := bac.getBasicAuthFromHeader("Bearer token")
		assert.Equal(t, ErrBasicNotSupported, err)
	})
}

func TestGetUserFromBasicAuth(t *testing.T) {
	bac := NewBasicAuthCollection()
	assert.Equal(t, "user", bac.getUserFromBasicAuth("dXNlcjpwYXNzMA=="))
	assert.Equal(t, "(invalid base64)", bac.getUserFromBasicAuth("invalid_base64"))
	assert.Equal(t, "(invalid basic auth format)", bac.getUserFromBasicAuth("dXNlcg=="))
}

func TestBasicAuthCollectionIsEmpty(t *testing.T) {
	t.Run("Is empty", func(t *testing.T) {
		bac := NewBasicAuthCollection()
		assert.True(t, bac.IsEmpty())
	})

	t.Run("Is not empty", func(t *testing.T) {
		bac := NewBasicAuthCollection()
		bac.Add("user", "password")
		assert.False(t, bac.IsEmpty())
	})
}

func TestBasicAuthCollectionImport(t *testing.T) {
	t.Run("Invalid JSON", func(t *testing.T) {
		bac := NewBasicAuthCollection()
		err := bac.Import("")
		assert.NotNil(t, err)
		assert.True(t, bac.IsEmpty())
	})
	t.Run("JSON is an empty map", func(t *testing.T) {
		bac := NewBasicAuthCollection()
		err := bac.Import("{}")
		assert.Nil(t, err)
		assert.True(t, bac.IsEmpty())
	})
	t.Run("JSON is non-empty map", func(t *testing.T) {
		bac := NewBasicAuthCollection()
		err := bac.Import(`{"user1": "password1", "user2": "password2"}`)
		assert.Nil(t, err)
		assert.False(t, bac.IsEmpty())
	})
}

func TestBasicAuthCollectionAllowedUserFromAuthorizationHeader(t *testing.T) {
	bac := NewBasicAuthCollection()
	bac.Add("user1", "password1")

	t.Run("Not a supported authorization header", func(t *testing.T) {
		_, err := bac.AllowedUserFromAuthorizationHeader("Bearer ABCDEF==")
		assert.Equal(t, ErrBasicNotSupported, err)
	})
	t.Run("Missing password", func(t *testing.T) {
		_, err := bac.AllowedUserFromAuthorizationHeader("BASIC dXNlcg==") // user
		assert.Equal(t, ErrBasicUnknownUser, err)
	})
	t.Run("Empty password", func(t *testing.T) {
		_, err := bac.AllowedUserFromAuthorizationHeader("basic dXNlcjo=") // user:
		assert.Equal(t, ErrBasicUnknownUser, err)
	})
	t.Run("Unknown user", func(t *testing.T) {
		_, err := bac.AllowedUserFromAuthorizationHeader("Basic dXNlcjpwYXNzd29yZA==") // user:password
		assert.Equal(t, ErrBasicUnknownUser, err)
	})
	t.Run("Unknown user without padding with =", func(t *testing.T) {
		_, err := bac.AllowedUserFromAuthorizationHeader("Basic dXNlcjpwYXNzd29yZA") // user:password
		assert.Equal(t, ErrBasicUnknownUser, err)
	})
	t.Run("Valid user", func(t *testing.T) {
		user, err := bac.AllowedUserFromAuthorizationHeader("Basic dXNlcjE6cGFzc3dvcmQx") // user1:password1
		assert.Nil(t, err)
		assert.Equal(t, "user1", *user)
	})
}
