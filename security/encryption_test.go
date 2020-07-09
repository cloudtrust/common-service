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
	assert.Nil(t, err)

	var res []byte
	res, err = encryption.Decrypt(encrypted, additional)
	assert.Nil(t, err)
	assert.Equal(t, value, res)
}

// Test our encryption in a normal encrypt/decrypt cycle
func TestAesGcmEncryptingFromBase64(t *testing.T) {
	var key = "ABCDEFGHIJKLMNOPQRSTUVWXYZ012345"
	var b64Key = base64.StdEncoding.EncodeToString([]byte(key))

	var encryption, err = NewAesGcmEncrypterFromBase64(b64Key, 16)
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
	var encryption, err = NewAesGcmEncrypter([]byte{0, 1, 2}, 16)
	assert.NotNil(t, err)

	// Generate a sample input (random=invalid... but won't be used)
	var input = make([]byte, 20)
	rand.Read(input)
	_, err = encryption.Decrypt(input, []byte("AesGcmInvalidKeySize"))
	assert.NotNil(t, err)
}

// Try to decrypt an invalid value
func TestAesGcmDecryptInvalidInput(t *testing.T) {
	// Generate key
	var key = make([]byte, 16)
	rand.Read(key)

	// Generate too short input
	var input = make([]byte, 10)
	rand.Read(input)

	var encryption, _ = NewAesGcmEncrypter(key, 16)
	_, err := encryption.Decrypt(input, []byte("AesGcmDecryptInvalidInput"))
	assert.NotNil(t, err)
}

// Try to decrypt with an invalid tag size
func TestAesGcmDecryptInvalidTagSize(t *testing.T) {
	// Generate key
	var key = make([]byte, 16)
	rand.Read(key)

	var encryption, _ = NewAesGcmEncrypter(key, 3)

	// Try to encrypt a value with a bad tag size
	_, err := encryption.Encrypt([]byte("Any value should not change the result of this test"), []byte("any additional data"))
	assert.NotNil(t, err)

	// Try to decrypt a value with a bad tag size
	_, err = encryption.Decrypt([]byte("Any value should not change the result of this test"), []byte("any additional data"))
	assert.NotNil(t, err)
}
