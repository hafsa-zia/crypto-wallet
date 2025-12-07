package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"math/big"
)

// PrivateFromHex: given D (private scalar) in hex, reconstruct priv+pub.
func PrivateFromHex(privHex string) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	dBytes, err := hex.DecodeString(privHex)
	if err != nil {
		return nil, nil, err
	}
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(dBytes)

	if d.Sign() <= 0 || d.Cmp(curve.Params().N) >= 0 {
		return nil, nil, errors.New("invalid private scalar D")
	}

	x, y := curve.ScalarBaseMult(d.Bytes())
	priv := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		},
		D: d,
	}
	return priv, &priv.PublicKey, nil
}
