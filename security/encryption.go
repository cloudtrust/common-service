package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"

	errorsMsg "github.com/cloudtrust/common-service/errors"
)

// EncrypterDecrypter used to encrypt/decrypt data
type EncrypterDecrypter interface {
	Encrypt(value []byte, additional []byte) ([]byte, error)
	Decrypt(value []byte, kid string, additional []byte) ([]byte, error)
	GetCurrentKeyID() string
}

type aesGcmKey struct {
	Kid      string `json:"kid"`
	Key      []byte `json:"value"`
	priority int    `json:"-"`
}

type keyMaterial struct {
	keys    []aesGcmKey
	tagSize int
}

// NewAesGcmEncrypterFromBase64 creation from json structure serialized as string
func NewAesGcmEncrypterFromBase64(keys string, tagSize int) (EncrypterDecrypter, error) {
	// parse key array
	var keyEntries []aesGcmKey
	err := json.Unmarshal([]byte(keys), &keyEntries)
	if err != nil {
		return nil, err
	}
	for i, k := range keyEntries {
		k.priority, err = strconv.Atoi(k.Kid[strings.LastIndex(k.Kid, "_")+1:])
		if err != nil {
			return nil, err
		}
		keyEntries[i] = k
	}
	// sort key entries according to priority
	sort.Slice(keyEntries, func(i, j int) bool {
		return keyEntries[i].priority > keyEntries[j].priority
	})
	km := keyMaterial{keys: keyEntries, tagSize: tagSize}
	// validate the correctness of the current key
	err = km.validate()
	if err != nil {
		return nil, err
	}
	return &km, nil
}

func (km *keyMaterial) GetCurrentKeyID() string {
	return km.keys[0].Kid
}

func (km *keyMaterial) validate() error {
	var sample = "This is a sample value used to check encrypt/decrypt cycle"
	var additionalData = []byte("additional data")

	for _, key := range km.keys {
		// create temporary key material for validation
		var kmSpecific = keyMaterial{
			keys:    []aesGcmKey{key},
			tagSize: km.tagSize,
		}
		var encrypted []byte
		kid := kmSpecific.GetCurrentKeyID()
		var err error
		if encrypted, err = kmSpecific.Encrypt([]byte(sample), additionalData); err != nil {
			return err
		}

		var decrypted []byte
		if decrypted, err = kmSpecific.Decrypt(encrypted, kid, additionalData); err != nil {
			return err
		}

		if string(decrypted) != sample {
			return errors.New(errorsMsg.MsgErrUnknown + "." + errorsMsg.EncryptDecrypt)
		}
	}
	return nil
}

func (km *keyMaterial) Encrypt(value []byte, additional []byte) ([]byte, error) {
	// select the most recent key
	key := km.keys[0]
	var block, err = aes.NewCipher(key.Key)
	if err != nil {
		return nil, err
	}

	var iv = make([]byte, 12)
	_, _ = rand.Read(iv)

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCMWithTagSize(block, km.tagSize)
	if err != nil {
		return nil, err
	}

	var enc = aesgcm.Seal(nil, iv, value, additional)
	encValue := append(iv, enc...)

	return encValue, err
}

func (km *keyMaterial) Decrypt(encData []byte, kid string, additional []byte) ([]byte, error) {
	// select the appropriate key
	var key aesGcmKey
	var found bool
	for _, k := range km.keys {
		if k.Kid == kid {
			key = k
			found = true
			break
		}
	}
	if !found {
		// key for decryption is not available
		return nil, errors.New(errorsMsg.MsgErrDecryptionKeyNotAvailable + "." + errorsMsg.EncryptDecrypt)
	}

	// decryption process
	if len(encData) <= 12 {
		return nil, errors.New(errorsMsg.MsgErrInvalidLength + "." + errorsMsg.Ciphertext)
	}

	var iv = encData[0:12]
	var encrypted = encData[12:]
	block, err := aes.NewCipher(key.Key)
	if err != nil {
		return nil, err
	}

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCMWithTagSize(block, km.tagSize)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, iv, encrypted, additional)
}
