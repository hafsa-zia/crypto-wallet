package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
)

func GenerateKeyPair() (string, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(key)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	pubBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", "", err
	}
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})

	return string(pubPem), string(privPem), nil
}

func getAESKey() ([]byte, error) {
	keyBytes, err := hex.DecodeString(config.AppConfig.AESSecretKey)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("AES key must be 32 bytes (64 hex chars)")
	}
	return keyBytes, nil
}

func EncryptPrivateKey(plainHex string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}

	plaintext, err := hex.DecodeString(plainHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func DecryptPrivateKey(encHex string) (string, error) {
	key, err := getAESKey()
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(encHex)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(plain), nil
}
