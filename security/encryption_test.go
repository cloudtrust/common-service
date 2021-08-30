package security

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testAesGcm(t *testing.T, encryption EncrypterDecrypter, value []byte) {
	var additional = make([]byte, 10)
	rand.Read(additional)
	var encrypted, err = encryption.Encrypt([]byte(value), additional)
	kid := encryption.GetCurrentKeyID()
	assert.Nil(t, err)

	var res []byte
	res, err = encryption.Decrypt(encrypted, kid, additional)
	assert.Nil(t, err)
	assert.Equal(t, value, res)
}

// Test our encryption in a normal encrypt/decrypt cycle
func TestAesGcmEncryptingFromBase64(t *testing.T) {
	var keys = `[
		{"kid":"DBB_1","value":"ABCDEFGHIJKLMNOPQRSTUVWXYZ012345"},
		{"kid":"DBB_2","value":"ABCDEFGHIJKLMNOPQRSTUVWXYZ012346"}
	]`

	var encryption, err = NewAesGcmEncrypterFromBase64(keys, 16)
	assert.Nil(t, err)
	assert.NotNil(t, encryption)

	testAesGcm(t, encryption, []byte("Sample value used in an encrypt/decrypt cycle to check our encryption tool"))

	var anotherSample = make([]byte, 300)
	rand.Read(anotherSample)
	testAesGcm(t, encryption, anotherSample)
}

// Try to test our encryption with an invalid base64 encoded key
func TestAesGcmFromBase64WithInvalidKey(t *testing.T) {
	var _, err = NewAesGcmEncrypterFromBase64("A", 16)
	assert.NotNil(t, err)
}

// Try to test our encryption with an invalid key. Valid keys should be 16, 24 or 32 bytes length
func TestAesGcmInvalidKeySize(t *testing.T) {
	var key = `[
		{"kid":"DBB_1","value":"aqNe"}
	]`
	var _, err = NewAesGcmEncrypterFromBase64(key, 16)
	assert.NotNil(t, err)
}

// Try to decrypt an invalid value
func TestAesGcmDecryptInvalidInput(t *testing.T) {
	// Generate key
	var key = make([]byte, 16)
	rand.Read(key)

	var keyMaterial = `[
		{"kid":"DBB_1","value":"` + base64.StdEncoding.EncodeToString(key) + `"}
	]`

	// Generate too short input
	var input = make([]byte, 10)
	rand.Read(input)

	var encryption, _ = NewAesGcmEncrypterFromBase64(keyMaterial, 16)
	_, err := encryption.Decrypt(input, "DBB_1", []byte("AesGcmDecryptInvalidInput"))
	assert.NotNil(t, err)
}

// Try to decrypt with an invalid tag size
func TestAesGcmDecryptInvalidTagSize(t *testing.T) {
	// Generate key
	var key = make([]byte, 16)
	rand.Read(key)

	var keyMaterial = `[
		{"kid":"DBB_1","value":"` + base64.StdEncoding.EncodeToString(key) + `"}
	]`

	var _, err = NewAesGcmEncrypterFromBase64(keyMaterial, 3)
	assert.NotNil(t, err)
}
