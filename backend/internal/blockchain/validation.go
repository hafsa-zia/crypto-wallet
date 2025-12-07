package blockchain

import (
	"errors"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
)

func ValidateBlock(b models.Block, prev *models.Block) error {
	if prev != nil {
		if b.PreviousHash != prev.Hash {
			return errors.New("invalid previous hash")
		}
		if b.Index != prev.Index+1 {
			return errors.New("invalid index")
		}
	}
	// check PoW, tx, etc...
	return nil
}
