package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/api"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
)

func main() {
	config.LoadConfig()
	if err := db.ConnectMongo(); err != nil {
		log.Fatal("failed to connect Mongo:", err)
	}

	r := gin.Default()

	// Add middlewares: logger + CORS
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORSMiddleware())

	api.RegisterRoutes(r)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
