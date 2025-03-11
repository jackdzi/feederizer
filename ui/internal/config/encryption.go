package config

import (
	"encoding/base64"
)

func Encrypt(plainText string) (string, error) {
	return base64.URLEncoding.EncodeToString([]byte(plainText)), nil
}

func Decrypt(cryptoText string) (string, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	return string(cipherText), nil
}
