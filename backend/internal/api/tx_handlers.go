package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appCrypto "github.com/hafsa-zia/crypto-wallet-backend/internal/crypto"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/utxo"
	"go.mongodb.org/mongo-driver/bson"
)

type CreateTxRequest struct {
	ReceiverWallet string  `json:"receiver_wallet" binding:"required"`
	Amount         float64 `json:"amount" binding:"required,gt=0"`
	Note           string  `json:"note"`
}

// in api/tx_handlers.go
func walletExists(ctx context.Context, walletID string) bool {
	err := db.Col("wallets").FindOne(ctx, bson.M{"wallet_id": walletID}).Err()
	return err == nil
}
func GetTxHistory(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	ctx := context.Background()

	col := db.Col("transactions")
	filter := bson.M{
		"$or": []bson.M{
			{"sender_wallet": walletID},
			{"receiver_wallet": walletID},
		},
	}
	cur, err := col.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer cur.Close(ctx)

	var txs []models.Transaction
	for cur.Next(ctx) {
		var t models.Transaction
		if err := cur.Decode(&t); err != nil {
			continue
		}
		txs = append(txs, t)
	}
	c.JSON(http.StatusOK, gin.H{"transactions": txs})
}

func CreateTransaction(c *gin.Context) {
	// user from JWT middleware
	walletID := c.GetString("wallet_id")

	var req CreateTxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if !walletExists(ctx, req.ReceiverWallet) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid receiver wallet"})
		return
	}

	// UTXO selection
	selected, change, err := utxo.SelectUTXOsForAmount(ctx, walletID, req.Amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		return
	}

	// build payload for signing
	timestamp := time.Now().UTC()
	payload := walletID + req.ReceiverWallet + fmt.Sprintf("%f", req.Amount) + timestamp.Format(time.RFC3339) + req.Note

	// load sender keys from DB
	userCol := db.Col("users")
	var user models.User
	if err := userCol.FindOne(ctx, bson.M{"wallet_id": walletID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	privHex, err := appCrypto.DecryptPrivateKey(user.EncryptedPrivKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("decrypt key failed: %v", err)})
		return
	}

	privKey, pubKey, err := appCrypto.PrivateFromHex(privHex)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("invalid private key: %v", err)})
		return
	}
	sig, err := appCrypto.SignMessage(privKey, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "sign failed"})
		return
	}

	// verify (for safety)
	ok, _ := appCrypto.VerifySignature(*pubKey, payload, sig)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "signature invalid"})
		return
	}

	inputs := []models.TxUTXOInput{}
	for _, u := range selected {
		inputs = append(inputs, models.TxUTXOInput{UTXOId: u.ID, Index: u.Index})
	}

	outputs := []models.TxUTXOOutput{
		{OwnerWallet: req.ReceiverWallet, Amount: req.Amount},
	}
	if change > 0 {
		outputs = append(outputs, models.TxUTXOOutput{OwnerWallet: walletID, Amount: change})
	}

	tx := models.Transaction{
		SenderWallet:   walletID,
		ReceiverWallet: req.ReceiverWallet,
		Amount:         req.Amount,
		Note:           req.Note,
		Timestamp:      timestamp,
		SenderPubKey:   user.PublicKey,
		Signature:      sig,
		Inputs:         inputs,
		Outputs:        outputs,
		Type:           "normal",
		Status:         "pending",
	}

	txCol := db.Col("pending_transactions")
	_, err = txCol.InsertOne(ctx, tx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction created (pending mining)"})
}
