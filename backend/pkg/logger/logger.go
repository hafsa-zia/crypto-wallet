package logger

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
)

func AddSystemLog(c *gin.Context, event, details string) {
	col := db.Col("system_logs")

	log := models.SystemLog{
		Event:     event,
		Details:   details,
		IP:        c.ClientIP(),
		Timestamp: time.Now().UTC(),
	}

	_, _ = col.InsertOne(context.Background(), log)
}
