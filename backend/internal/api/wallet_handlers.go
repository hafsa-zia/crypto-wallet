package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/utxo"
	"go.mongodb.org/mongo-driver/bson"
)

func GetWalletProfile(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	ctx := context.Background()

	usersCol := db.Col("users")
	var user models.User
	if err := usersCol.FindOne(ctx, bson.M{"wallet_id": walletID}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	balance, err := utxo.GetBalance(ctx, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "balance error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"full_name":      user.FullName,
		"email":          user.Email,
		"cnic":           user.CNIC,
		"wallet_id":      user.WalletID,
		"beneficiaries":  user.Beneficiaries,
		"zakat_deducted": user.ZakatDeducted,
		"balance":        balance,
	})
}

func GetBalance(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	ctx := context.Background()

	balance, err := utxo.GetBalance(ctx, walletID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "balance error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

func GetUTXOs(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	ctx := context.Background()

	col := db.Col("utxos")
	cur, err := col.Find(ctx, bson.M{"owner_wallet": walletID, "is_spent": false})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer cur.Close(ctx)

	var out []models.UTXO
	for cur.Next(ctx) {
		var u models.UTXO
		if err := cur.Decode(&u); err != nil {
			continue
		}
		out = append(out, u)
	}
	c.JSON(http.StatusOK, gin.H{"utxos": out})
}

type BeneficiariesReq struct {
	Beneficiaries []string `json:"beneficiaries"`
}

func UpdateBeneficiaries(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	var req BeneficiariesReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	_, err := db.Col("users").UpdateOne(ctx,
		bson.M{"wallet_id": walletID},
		bson.M{"$set": bson.M{"beneficiaries": req.Beneficiaries}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "beneficiaries updated"})
}
