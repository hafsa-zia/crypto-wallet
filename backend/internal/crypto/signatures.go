package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/hex"
	"errors"
	"math/big"
)

type ecdsaSignature struct {
	R, S *big.Int
}

func SignMessage(priv *ecdsa.PrivateKey, msg string) (string, error) {
	hash := sha256.Sum256([]byte(msg))
	r, s, err := ecdsa.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return "", err
	}

	sigBytes, err := asn1.Marshal(ecdsaSignature{R: r, S: s})
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(sigBytes), nil
}

func VerifySignature(pub ecdsa.PublicKey, msg, sigHex string) (bool, error) {
	hash := sha256.Sum256([]byte(msg))

	sigBytes, err := hex.DecodeString(sigHex)
	if err != nil {
		return false, err
	}

	var sig ecdsaSignature
	_, err = asn1.Unmarshal(sigBytes, &sig)
	if err != nil {
		return false, err
	}

	ok := ecdsa.Verify(&pub, hash[:], sig.R, sig.S)
	if !ok {
		return false, errors.New("invalid signature")
	}
	return true, nil
}
