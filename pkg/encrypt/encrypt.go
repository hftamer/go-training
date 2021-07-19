package encrypt

// The main idea is taken from: https://www.thepolyglotdeveloper.com/2018/02/encrypt-decrypt-data-golang-application-crypto-packages/

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

func Decrypt(passphrase string, data string) (string, error) {
	if data == "" {
		return "", nil
	}
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed creating aes block, %s", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed creating gcm, %s", err)
	}
	nonceSize := gcm.NonceSize()
	dataByte := []byte(data)
	nonce, ciphertext := dataByte[:nonceSize], dataByte[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed getting plaintext, %s", err)
	}
	return string(plaintext), nil
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(passphrase string, data string) (string, error) {
	dataByte := []byte(data)
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, dataByte, nil)
	return string(ciphertext), nil
}
