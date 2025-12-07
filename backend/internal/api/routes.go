package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth
	api.POST("/auth/request-otp", RequestOTPHandler) // ✅ new
	api.POST("/auth/register", Register)
	api.POST("/auth/login", Login)

	// Protected routes
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())

	// Wallet
	protected.GET("/wallet", GetWalletProfile)
	protected.GET("/wallet/balance", GetBalance)
	protected.GET("/wallet/utxos", GetUTXOs)
	protected.POST("/wallet/beneficiaries", UpdateBeneficiaries)

	// Transactions
	protected.POST("/tx", CreateTransaction)
	protected.GET("/tx/history", GetTxHistory)

	// Blockchain
	protected.GET("/blocks", GetBlocks)
	protected.GET("/blocks/:id", GetBlockByID)
	protected.POST("/admin/mine", MinePending)

	// Zakat
	protected.POST("/admin/run-zakat", RunZakatNowHandler)
	protected.GET("/reports/summary", GetReportsSummary)

	// Logs
	protected.GET("/logs/system", GetSystemLogs)
	protected.GET("/logs/transactions", GetTxLogs)
}

// RunZakatNowHandler handles POST /admin/run-zakat requests.
func RunZakatNowHandler(c *gin.Context) {
	// TODO: Implement logic to run zakat now
	c.JSON(http.StatusOK, gin.H{"message": "Zakat run initiated"})
}

// GetBlockByID handles GET /blocks/:id requests.
func GetBlockByID(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to fetch block by ID
	c.JSON(http.StatusOK, gin.H{"message": "Block details for ID " + id})
}
