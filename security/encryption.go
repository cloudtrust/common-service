package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"

	errorsMsg "github.com/cloudtrust/common-service/errors"
)

// CrypterDecrypter used to encrypt/decrypt data
type CrypterDecrypter interface {
	Encrypt(value []byte, additional []byte) ([]byte, error)
	Decrypt(value []byte, additional []byte) ([]byte, error)
}

type aesGcmCrypting struct {
	key     []byte
	tagSize int
}

// NewAesGcmEncrypter creation from slice of bytes
func NewAesGcmEncrypter(key []byte, tagSize int) (CrypterDecrypter, error) {
	cd := aesGcmCrypting{
		key:     key,
		tagSize: tagSize,
	}
	return &cd, cd.validate()
}

// NewAesGcmEncrypterFromBase64 creation from base64
func NewAesGcmEncrypterFromBase64(base64Key string, tagSize int) (CrypterDecrypter, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	return NewAesGcmEncrypter(key, tagSize)
}

func (cd *aesGcmCrypting) validate() error {
	var sample = "This is a sample value used to check encrypt/decrypt cycle"
	var additionalData = []byte("additional data")

	var encrypted []byte
	var err error
	if encrypted, err = cd.Encrypt([]byte(sample), additionalData); err != nil {
		return err
	}

	var decrypted []byte
	if decrypted, err = cd.Decrypt(encrypted, additionalData); err != nil {
		return err
	}

	if string(decrypted) == sample {
		return nil
	}
	return errors.New(errorsMsg.MsgErrUnknown + "." + errorsMsg.EncryptDecrypt)
}

func (cd *aesGcmCrypting) Encrypt(value []byte, additional []byte) ([]byte, error) {
	var block, err = aes.NewCipher(cd.key)
	if err != nil {
		return nil, err
	}

	var iv = make([]byte, 12)
	rand.Read(iv)

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCMWithTagSize(block, cd.tagSize)
	if err != nil {
		return nil, err
	}

	var enc = aesgcm.Seal(nil, iv, value, additional)

	return append(iv, enc...), err
}

func (cd *aesGcmCrypting) Decrypt(value []byte, additional []byte) ([]byte, error) {
	if len(value) <= 12 {
		return nil, errors.New(errorsMsg.MsgErrInvalidLength + "." + errorsMsg.Ciphertext)
	}

	var iv = value[0:12]
	var crypted = value[12:]
	var block, err = aes.NewCipher(cd.key)
	if err != nil {
		return nil, err
	}

	var aesgcm cipher.AEAD
	aesgcm, err = cipher.NewGCMWithTagSize(block, cd.tagSize)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, iv, crypted, additional)
}
