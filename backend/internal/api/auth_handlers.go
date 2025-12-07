package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appCrypto "github.com/hafsa-zia/crypto-wallet-backend/internal/crypto"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	CNIC     string `json:"cnic" binding:"required"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// POST /api/auth/register
func Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("bad_request email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	usersCol := db.Col("users")
	walletsCol := db.Col("wallets")

	// check if email exists
	var existing models.User
	err := usersCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&existing)
	if err == nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("email_exists email=%s", req.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "email already registered"})
		return
	}

	// hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("hash_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// generate keypair
	pubKeyPEM, privKeyPEM, err := appCrypto.GenerateKeyPair()
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("keypair_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate keypair"})
		return
	}

	// encrypt private key (simple base64 in our helper)
	encryptedPriv, err := appCrypto.EncryptPrivateKey(privKeyPEM)
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("encrypt_priv_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt private key"})
		return
	}

	// wallet id = SHA256(publicKeyPEM)
	h := sha256.Sum256([]byte(pubKeyPEM))
	walletID := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	user := models.User{
		ID:               primitive.NewObjectID().Hex(),
		FullName:         req.FullName,
		Email:            req.Email,
		CNIC:             req.CNIC,
		WalletID:         walletID,
		PasswordHash:     string(hashed),
		PublicKey:        pubKeyPEM,
		EncryptedPrivKey: encryptedPriv,
		CreatedAt:        now,
	}

	_, err = usersCol.InsertOne(ctx, user)
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("db_insert_user_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	// create wallet doc
	walletDoc := models.Wallet{
		ID:       primitive.NewObjectID().Hex(),
		UserID:   user.ID,
		WalletID: walletID,
		// CreatedAt field removed or replace with the correct field if needed
	}
	_, _ = walletsCol.InsertOne(ctx, walletDoc)

	// generate JWT
	token, err := middleware.GenerateToken(user.ID, user.WalletID)
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("token_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.AddSystemLog(c,
		"register_success",
		fmt.Sprintf("email=%s wallet=%s", user.Email, user.WalletID),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "registration successful",
		"token":   token,
		"user": gin.H{
			"full_name": user.FullName,
			"email":     user.Email,
			"wallet_id": user.WalletID,
		},
	})
}

// POST /api/auth/login
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.AddSystemLog(c,
			"login_failed",
			fmt.Sprintf("bad_request email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	usersCol := db.Col("users")

	var user models.User
	if err := usersCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&user); err != nil {
		logger.AddSystemLog(c,
			"login_failed",
			fmt.Sprintf("user_not_found email=%s", req.Email),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.AddSystemLog(c,
			"login_failed",
			fmt.Sprintf("wrong_password email=%s", req.Email),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.WalletID)
	if err != nil {
		logger.AddSystemLog(c,
			"login_failed",
			fmt.Sprintf("token_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	logger.AddSystemLog(c,
		"login_success",
		fmt.Sprintf("email=%s wallet=%s", user.Email, user.WalletID),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"full_name": user.FullName,
			"email":     user.Email,
			"wallet_id": user.WalletID,
		},
	})
}
