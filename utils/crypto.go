package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var secretKey = []byte(os.Getenv("AES_SECRET_KEY"))

func Decrypt(encryptedText string) (string, error) {

	data, err := base64.RawURLEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()

	nonce, cipherText := data[:nonceSize], data[nonceSize:]

	plainText, err := aesGCM.Open(
		nil,
		nonce,
		cipherText,
		nil,
	)

	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
