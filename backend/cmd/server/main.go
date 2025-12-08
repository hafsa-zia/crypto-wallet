package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/api"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
)

func main() {
	// Load config
	config.LoadConfig()

	// Connect to MongoDB
	if err := db.ConnectMongo(); err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Create Gin router
	r := gin.Default()

	// Middlewares
	r.Use(middleware.RequestLogger())
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Backend is running"})
	})

	// Register API routes
	api.RegisterRoutes(r)

	// Start server
	port := ":8080"
	log.Printf("🚀 Backend is running on port %s", port)
	if err := r.Run(port); err != nil {
		log.Fatal(err)
	}
}
