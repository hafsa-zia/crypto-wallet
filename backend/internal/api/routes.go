package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/zakat"
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
	// Profile
	protected.GET("/profile", GetProfile)
	protected.PUT("/profile", UpdateProfile)

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
	protected.POST("/zakat/run-self", RunSelfZakatHandler)
	protected.GET("/reports/summary", GetReportsSummary)

	// Logs
	protected.GET("/logs/system", GetSystemLogs)
	protected.GET("/logs/transactions", GetTxLogs)
}

// RunZakatNowHandler handles POST /admin/run-zakat requests.
// RunZakatNowHandler handles POST /api/admin/run-zakat
// It creates pending zakat_deduction transactions for all wallets
// and returns how many were created. To confirm them on-chain,
// you then mine pending txs via /api/admin/mine.
// RunSelfZakatHandler handles POST /api/zakat/run-self.
// It runs zakat (2.5%) for the CURRENT logged-in wallet and creates
// a "zakat_deduction" transaction with status "pending".
func RunSelfZakatHandler(c *gin.Context) {
	walletID := c.GetString("wallet_id")
	if walletID == "" {
		if alt := c.GetString("walletID"); alt != "" {
			walletID = alt
		} else if alt2 := c.GetString("wallet"); alt2 != "" {
			walletID = alt2
		}
	}

	if walletID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "wallet id not found in token context",
		})
		return
	}

	txDoc, err := zakat.RunZakatForWallet(c, walletID)
	if err != nil {
		// soft errors → 200
		if err.Error() == "no UTXOs / no balance, zakat not due" ||
			err.Error() == "balance <= 0, zakat not due" ||
			err.Error() == "zakat computed as 0, nothing to do" ||
			err.Error() == "not enough UTXOs to cover zakat" {
			c.JSON(http.StatusOK, gin.H{
				"message": err.Error(),
			})
			return
		}

		// real errors → 500
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to run zakat",
			"details": err.Error(),
		})
		return
	}

	if txDoc == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "No zakat due for this wallet.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Zakat transaction created as pending. Mine pending transactions to confirm.",
		"tx":      txDoc,
	})
}

// GetBlockByID handles GET /blocks/:id requests.
func GetBlockByID(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement logic to fetch block by ID
	c.JSON(http.StatusOK, gin.H{"message": "Block details for ID " + id})
}
