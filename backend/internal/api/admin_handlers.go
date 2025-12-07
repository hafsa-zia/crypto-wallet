package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FundRequest struct {
	WalletID string  `json:"wallet_id" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
}

// POST /api/admin/fund
// Simple faucet: create a UTXO for the given wallet
func AdminFundWallet(c *gin.Context) {
	var req FundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()

	// check wallet exists
	walletCol := db.Col("wallets")
	var w models.Wallet
	if err := walletCol.FindOne(ctx, bson.M{"wallet_id": req.WalletID}).Decode(&w); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet does not exist"})
		return
	}

	// create a synthetic tx id for faucet
	txID := "faucet-" + primitive.NewObjectID().Hex()

	utxoCol := db.Col("utxos")
	_, err := utxoCol.InsertOne(ctx, models.UTXO{
		TxID:        txID,
		Index:       0,
		OwnerWallet: req.WalletID,
		Amount:      req.Amount,
		IsSpent:     false,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert UTXO"})
		return
	}

	// optional: log a transaction record of type "faucet"
	txCol := db.Col("transactions")
	_, _ = txCol.InsertOne(ctx, models.Transaction{
		ID:             txID,
		SenderWallet:   "SYSTEM_FAUCET",
		ReceiverWallet: req.WalletID,
		Amount:         req.Amount,
		Note:           "Admin faucet funding",
		Type:           "faucet",
		Status:         "confirmed",
		Timestamp:      time.Now().UTC(),
	})

	c.JSON(http.StatusOK, gin.H{
		"message":   "wallet funded via faucet",
		"wallet_id": req.WalletID,
		"amount":    req.Amount,
	})
}
