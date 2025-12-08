package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"io"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
)

// GenerateKeyPair generates an ECDSA P-256 keypair and returns
// the public key (hex of X||Y) and private scalar D as hex.
func GenerateKeyPair() (string, string, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), crand.Reader)
	if err != nil {
		return "", "", err
	}

	// private scalar D as hex
	privHex := hex.EncodeToString(priv.D.Bytes())

	// public key bytes: X||Y
	xBytes := priv.PublicKey.X.Bytes()
	yBytes := priv.PublicKey.Y.Bytes()
	pubBytes := append(xBytes, yBytes...)
	pubHex := hex.EncodeToString(pubBytes)

	return pubHex, privHex, nil
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

// EncryptPrivateKey expects the private key as a hex string (plaintext hex),
// decodes it and encrypts the raw bytes with AES-GCM. It returns the ciphertext
// as a hex string (nonce + ciphertext).
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
	if _, err := io.ReadFull(crand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

// DecryptPrivateKey takes the hex-encoded ciphertext (nonce + ciphertext),
// decrypts it using AES-GCM and returns the original private key as hex string.
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

	// Return the private key as hex string so callers (PrivateFromHex)
	// can reconstruct the ECDSA private key.
	return hex.EncodeToString(plain), nil
}
