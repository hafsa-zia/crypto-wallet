package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"go.mongodb.org/mongo-driver/bson"
)

type reportTotals struct {
	Total float64 `bson:"total" json:"total"`
	Count int64   `bson:"count" json:"count"`
}

// GET /api/reports/summary
func GetReportsSummary(c *gin.Context) {
	ctx := context.Background()
	walletID := c.GetString("wallet_id")
	if walletID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing wallet id"})
		return
	}

	txCol := db.Col("transactions")

	// ---- total sent ----
	sentAgg, _ := txCol.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"sender_wallet": walletID, "status": "confirmed"}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
			"count": bson.M{"$sum": 1},
		}},
	})
	var sent reportTotals
	if sentAgg.Next(ctx) {
		_ = sentAgg.Decode(&sent)
	}
	sentAgg.Close(ctx)

	// ---- total received ----
	recvAgg, _ := txCol.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"receiver_wallet": walletID, "status": "confirmed"}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
			"count": bson.M{"$sum": 1},
		}},
	})
	var recv reportTotals
	if recvAgg.Next(ctx) {
		_ = recvAgg.Decode(&recv)
	}
	recvAgg.Close(ctx)

	// ---- zakat deducted (by type) ----
	zakatAgg, _ := txCol.Aggregate(ctx, []bson.M{
		{"$match": bson.M{"sender_wallet": walletID, "type": "zakat_deduction"}},
		{"$group": bson.M{
			"_id":   nil,
			"total": bson.M{"$sum": "$amount"},
			"count": bson.M{"$sum": 1},
		}},
	})
	var zakat reportTotals
	if zakatAgg.Next(ctx) {
		_ = zakatAgg.Decode(&zakat)
	}
	zakatAgg.Close(ctx)

	c.JSON(http.StatusOK, gin.H{
		"total_sent_amount":       sent.Total,
		"total_sent_count":        sent.Count,
		"total_received_amount":   recv.Total,
		"total_received_count":    recv.Count,
		"zakat_deducted_amount":   zakat.Total,
		"zakat_deducted_tx_count": zakat.Count,
	})
}
