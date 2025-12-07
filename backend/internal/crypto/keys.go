package crypto

import (
	"crypto/ecdsa"

	"crypto/sha256"
	"encoding/hex"
)

type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
}

// Removed duplicate GenerateKeyPair function to resolve redeclaration error.

func PublicKeyToBytes(pub ecdsa.PublicKey) []byte {
	xBytes := pub.X.Bytes()
	yBytes := pub.Y.Bytes()
	return append(xBytes, yBytes...)
}

func WalletIDFromPublicKey(pub ecdsa.PublicKey) string {
	h := sha256.Sum256(PublicKeyToBytes(pub))
	return hex.EncodeToString(h[:])
}
