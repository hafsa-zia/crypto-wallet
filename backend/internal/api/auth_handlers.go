package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	appCrypto "github.com/hafsa-zia/crypto-wallet-backend/internal/crypto"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/email"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/middleware"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	CNIC     string `json:"cnic" binding:"required"`
	OTP      string `json:"otp" binding:"required"`
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
	otpCol := db.Col("email_verifications")

	// ✅ OTP validation BEFORE we create the user
	var otpDoc models.EmailOTP
	if err := otpCol.FindOne(ctx, bson.M{"email": req.Email}).Decode(&otpDoc); err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("otp_not_found email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No OTP found for this email. Please request a new OTP."})
		return
	}

	// Check if OTP is expired first
	if time.Now().After(otpDoc.ExpiresAt) {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("otp_expired email=%s expires_at=%v now=%v", req.Email, otpDoc.ExpiresAt, time.Now()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired. Please request a new OTP."})
		return
	}

	// Check if OTP code matches (case-sensitive)
	if otpDoc.OTP != req.OTP {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("otp_invalid email=%s provided=%s stored=%s", req.Email, req.OTP, otpDoc.OTP),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP"})
		return
	}

	// Check if OTP was already used
	if otpDoc.Verified {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("otp_already_used email=%s", req.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "OTP already used for this email. Please request a new OTP."})
		return
	}

	// mark OTP as used (non-fatal if it fails)
	_, _ = otpCol.UpdateOne(
		ctx,
		bson.M{"_id": otpDoc.ID},
		bson.M{"$set": bson.M{"verified": true}},
	)

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

	// generate keypair (returns both as hex strings)
	pubKeyHex, privKeyHex, err := appCrypto.GenerateKeyPair()
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("keypair_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate keypair"})
		return
	}

	// encrypt private key (already hex, pass directly)
	encryptedPriv, err := appCrypto.EncryptPrivateKey(privKeyHex)
	if err != nil {
		logger.AddSystemLog(c,
			"register_failed",
			fmt.Sprintf("encrypt_priv_error email=%s error=%s", req.Email, err.Error()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to encrypt private key"})
		return
	}

	// wallet id = SHA256(pubKeyHex)
	h := sha256.Sum256([]byte(pubKeyHex))
	walletID := hex.EncodeToString(h[:])

	now := time.Now().UTC()
	user := models.User{
		ID:               primitive.NewObjectID().Hex(),
		FullName:         req.FullName,
		Email:            req.Email,
		CNIC:             req.CNIC,
		WalletID:         walletID,
		PasswordHash:     string(hashed),
		PublicKey:        pubKeyHex,
		EncryptedPrivKey: encryptedPriv,
		Beneficiaries:    []string{},
		ZakatDeducted:    0,
		CreatedAt:        now,
		UpdatedAt:        now,
		EmailVerified:    true, // ✅ OTP passed
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
			"full_name":      user.FullName,
			"email":          user.Email,
			"wallet_id":      user.WalletID,
			"email_verified": user.EmailVerified,
		},
	})
}

// RequestOTPHandler handles POST /api/auth/request-otp
// Body: { "email": "user@example.com" }
func RequestOTPHandler(c *gin.Context) {
	type reqBody struct {
		Email string `json:"email" binding:"required,email"`
	}

	var body reqBody
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Valid email is required"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Generate 6-digit OTP
	rand.Seed(time.Now().UnixNano())
	otp := fmt.Sprintf("%06d", rand.Intn(1000000))

	coll := db.Col("email_verifications")

	// Upsert OTP doc for this email
	_, err := coll.UpdateOne(
		ctx,
		bson.M{"email": body.Email},
		bson.M{
			"$set": bson.M{
				"email":      body.Email,
				"otp":        otp,
				"expires_at": time.Now().Add(10 * time.Minute),
				"verified":   false,
			},
		},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store OTP"})
		return
	}

	if err := email.SendOTPEmail(body.Email, otp); err != nil {
		// now we TELL the frontend about the failure
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to send OTP email: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent to your email address"})

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
			"full_name":      user.FullName,
			"email":          user.Email,
			"wallet_id":      user.WalletID,
			"email_verified": user.EmailVerified,
		},
	})
}
