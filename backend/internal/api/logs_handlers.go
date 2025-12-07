package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GET /api/logs/system
func GetSystemLogs(c *gin.Context) {
	ctx := context.Background()
	col := db.Col("system_logs")

	cur, err := col.Find(ctx, bson.M{}, options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(100))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer cur.Close(ctx)

	var logs []models.SystemLog
	for cur.Next(ctx) {
		var l models.SystemLog
		if err := cur.Decode(&l); err != nil {
			continue
		}
		logs = append(logs, l)
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs})
}

// GET /api/logs/transactions
// Optional: last 100 transactions involving current wallet (acts as tx log)
func GetTxLogs(c *gin.Context) {
	ctx := context.Background()
	walletID := c.GetString("wallet_id")
	if walletID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing wallet id"})
		return
	}

	col := db.Col("transactions")

	filter := bson.M{
		"$or": []bson.M{
			{"sender_wallet": walletID},
			{"receiver_wallet": walletID},
		},
	}

	cur, err := col.Find(ctx, filter, options.Find().
		SetSort(bson.M{"timestamp": -1}).
		SetLimit(100))
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

	c.JSON(http.StatusOK, gin.H{"logs": txs})
}
