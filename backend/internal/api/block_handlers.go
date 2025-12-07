package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/blockchain"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models" // 👈 here
	"github.com/hafsa-zia/crypto-wallet-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const miningRewardAmount = 50.0 // reward per mined block

// GET /api/blocks
func GetBlocks(c *gin.Context) {
	ctx := context.Background()
	cur, err := db.Col("blocks").Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"index": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer cur.Close(ctx)

	var blocks []models.Block
	for cur.Next(ctx) {
		var b models.Block
		if err := cur.Decode(&b); err != nil {
			continue
		}
		blocks = append(blocks, b)
	}
	c.JSON(http.StatusOK, gin.H{"blocks": blocks})
}

// POST /api/admin/mine
// Mines a block with: 1) mining reward, 2) all pending user transactions.
func MinePending(c *gin.Context) {
	ctx := context.Background()

	minerWalletID := c.GetString("wallet_id")
	if minerWalletID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing miner wallet in token"})
		return
	}

	blocksCol := db.Col("blocks")
	pendingCol := db.Col("pending_transactions")
	utxoCol := db.Col("utxos")
	txCol := db.Col("transactions")

	// --- get last block (for index + prev hash) ---
	var last models.Block
	err := blocksCol.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"index": -1})).Decode(&last)
	hasPrev := err == nil

	// --- load all pending user transactions ---
	cur, err := pendingCol.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error loading pending txs"})
		return
	}
	defer cur.Close(ctx)

	var pendingTxs []models.Transaction
	for cur.Next(ctx) {
		var t models.Transaction
		if err := cur.Decode(&t); err != nil {
			continue
		}
		pendingTxs = append(pendingTxs, t)
	}

	// --- create coinbase (mining reward) transaction ---
	coinbaseID := primitive.NewObjectID().Hex()
	coinbaseTx := models.Transaction{
		ID:             coinbaseID,
		SenderWallet:   "SYSTEM_COINBASE",
		ReceiverWallet: minerWalletID,
		Amount:         miningRewardAmount,
		Note:           "Mining reward",
		Type:           "mining_reward",
		Status:         "confirmed",
		Timestamp:      time.Now().UTC(),
	}

	// --- assemble all block transactions: reward + pending ---
	allTxs := make([]models.Transaction, 0, len(pendingTxs)+1)
	allTxs = append(allTxs, coinbaseTx)
	allTxs = append(allTxs, pendingTxs...)

	// --- build block ---
	block := models.Block{
		Index:        0,
		Timestamp:    time.Now().UTC(),
		Transactions: allTxs,
		PreviousHash: "",
	}
	if hasPrev {
		block.Index = last.Index + 1
		block.PreviousHash = last.Hash
	}

	// merkle + PoW
	block.MerkleRoot = blockchain.MerkleRoot(block.Transactions)
	blockchain.MineBlock(&block)

	// --- insert block ---
	_, err = blocksCol.InsertOne(ctx, block)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "block insert error"})
		return
	}

	// --- 1) insert reward UTXO & tx ---
	_, err = utxoCol.InsertOne(ctx, models.UTXO{
		TxID:        coinbaseTx.ID,
		OwnerWallet: minerWalletID,
		Amount:      miningRewardAmount,
		IsSpent:     false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert reward UTXO"})
		return
	}
	// save reward tx
	coinbaseTx.Status = "confirmed"
	_, _ = txCol.InsertOne(ctx, coinbaseTx)

	// --- 2) process each pending transaction: UTXO spend + create ---
	for _, t := range pendingTxs {
		// 2.1 collect sender's available UTXOs
		var senderUtxos []models.UTXO
		uCur, err := utxoCol.Find(ctx, bson.M{
			"owner_wallet": t.SenderWallet,
			"is_spent":     false,
		})
		if err != nil {
			continue
		}
		for uCur.Next(ctx) {
			var u models.UTXO
			if err := uCur.Decode(&u); err != nil {
				continue
			}
			senderUtxos = append(senderUtxos, u)
		}
		uCur.Close(ctx)

		// 2.2 pick enough UTXOs to cover tx.Amount
		var used []models.UTXO
		var total float64
		for _, u := range senderUtxos {
			used = append(used, u)
			total += u.Amount
			if total >= t.Amount {
				break
			}
		}

		// not enough balance at mining time? skip this tx
		if total < t.Amount {
			continue
		}

		// 2.3 mark used UTXOs as spent
		for _, u := range used {
			_, _ = utxoCol.UpdateOne(ctx,
				bson.M{
					"tx_id":        u.TxID,
					"owner_wallet": u.OwnerWallet,
					"amount":       u.Amount,
					"is_spent":     false,
				},
				bson.M{"$set": bson.M{"is_spent": true}},
			)
		}

		// 2.4 create receiver UTXO
		_, _ = utxoCol.InsertOne(ctx, models.UTXO{
			TxID:        t.ID,
			OwnerWallet: t.ReceiverWallet,
			Amount:      t.Amount,
			IsSpent:     false,
		})

		// 2.5 change UTXO back to sender if any leftover
		change := total - t.Amount
		if change > 0 {
			_, _ = utxoCol.InsertOne(ctx, models.UTXO{
				TxID:        t.ID,
				OwnerWallet: t.SenderWallet,
				Amount:      change,
				IsSpent:     false,
			})
		}

		// 2.6 mark tx as confirmed and save to main transactions collection
		t.Status = "confirmed"
		_, _ = txCol.InsertOne(ctx, t)
	}

	// 2.7 clear pending pool (since we've processed what's possible)
	_, _ = pendingCol.DeleteMany(ctx, bson.M{})
	logger.AddSystemLog(
		c,
		"mined_block",
		fmt.Sprintf("Block #%d mined by %s with reward %.4f, user_tx=%d",
			block.Index,
			minerWalletID,
			miningRewardAmount,
			len(pendingTxs),
		),
	)

	c.JSON(http.StatusOK, gin.H{
		"message":       "block mined with reward and pending transactions",
		"block_index":   block.Index,
		"block_hash":    block.Hash,
		"miner_wallet":  minerWalletID,
		"reward_amount": miningRewardAmount,
		"tx_in_block":   len(allTxs),
		"user_tx_mined": len(pendingTxs),
	})
}
