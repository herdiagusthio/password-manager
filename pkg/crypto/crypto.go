package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts plainText using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func Encrypt(plainText, key string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("crypto: invalid key length (must be 16, 24, or 32 bytes)")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Seal encrypts and authenticates plainText, appending the result to nonce.
	// The nonce is prepended to the ciphertext to be used for decryption.
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts cipherText (base64 encoded) using AES-GCM with the provided key.
func Decrypt(cipherText, key string) (string, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("crypto: invalid key length (must be 16, 24, or 32 bytes)")
	}

	data, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("crypto: ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
