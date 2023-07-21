package crypt

import (
	"crypto/aes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

func AesEncrypt(data []byte, key []byte) (string, error) {
	a, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	en := make([]byte, len(data))
	a.Encrypt(en, data)
	return base64.StdEncoding.EncodeToString(en), nil
}

func AesDecrypt(data string, key []byte) (string, error) {
	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	a, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	de := make([]byte, len(data))
	a.Decrypt(de, d)
	return string(de), nil
}

func GenShaHash(data []byte, length int) string {
	d := sha256.Sum256(data)
	return hex.EncodeToString(d[:length])
}
