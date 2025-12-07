package blockchain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
)

// MerkleRoot computes a simple merkle root for a slice of transactions.
func MerkleRoot(txs []models.Transaction) string {
	if len(txs) == 0 {
		return ""
	}

	var hashes [][]byte
	for _, tx := range txs {
		h := sha256.Sum256([]byte(tx.ID + tx.SenderWallet + tx.ReceiverWallet))
		hashes = append(hashes, h[:])
	}

	for len(hashes) > 1 {
		var next [][]byte
		for i := 0; i < len(hashes); i += 2 {
			if i+1 == len(hashes) {
				// odd number of nodes, duplicate last
				next = append(next, hashes[i])
			} else {
				concat := append(hashes[i], hashes[i+1]...)
				h := sha256.Sum256(concat)
				next = append(next, h[:])
			}
		}
		hashes = next
	}

	return hex.EncodeToString(hashes[0])
}
