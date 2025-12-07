package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
)

func calculateBlockHash(b models.Block) string {
	data := fmt.Sprintf("%d|%s|%s|%d|%s",
		b.Index,
		b.Timestamp.UTC().String(),
		b.PreviousHash,
		b.Nonce,
		b.MerkleRoot,
	)
	h := sha256.Sum256([]byte(data))
	return hex.EncodeToString(h[:])
}

func hasRequiredDifficulty(hash string) bool {
	prefix := strings.Repeat("0", config.AppConfig.PowDifficulty)
	return strings.HasPrefix(hash, prefix)
}

func MineBlock(b *models.Block) {
	nonce := int64(0)
	for {
		b.Nonce = nonce
		h := calculateBlockHash(*b)
		if hasRequiredDifficulty(h) {
			b.Hash = h
			return
		}
		nonce++
	}
}
