package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
)

// GET /api/profile
// Returns the currently logged-in user's profile.
func GetProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	ctx := context.Background()
	usersCol := db.Col("users")

	var user models.User
	if err := usersCol.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.ID,
		"full_name":  user.FullName,
		"email":      user.Email,
		"cnic":       user.CNIC,
		"wallet_id":  user.WalletID,
		"created_at": user.CreatedAt,
	})
}

type profileUpdateRequest struct {
	FullName *string `json:"full_name,omitempty"`
	CNIC     *string `json:"cnic,omitempty"`

	// For email change, both must be present:
	Email    *string `json:"email,omitempty"`
	EmailOTP *string `json:"email_otp,omitempty"`
}

// PUT /api/profile
// Allows updating full_name & cnic directly.
// Email change requires OTP verification (via email_verifications collection).
func UpdateProfile(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user id"})
		return
	}

	var req profileUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	ctx := context.Background()
	usersCol := db.Col("users")
	otpCol := db.Col("email_verifications")

	// Load current user
	var user models.User
	if err := usersCol.FindOne(ctx, bson.M{"_id": userID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

	updateFields := bson.M{}
	now := time.Now().UTC()

	// Name update
	if req.FullName != nil && *req.FullName != "" && *req.FullName != user.FullName {
		updateFields["full_name"] = *req.FullName
	}

	// CNIC update
	if req.CNIC != nil && *req.CNIC != "" && *req.CNIC != user.CNIC {
		updateFields["cnic"] = *req.CNIC
	}

	// Email update (requires OTP)
	if req.Email != nil && *req.Email != "" && *req.Email != user.Email {
		// Must have OTP
		if req.EmailOTP == nil || *req.EmailOTP == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP is required to change email"})
			return
		}

		type emailVerificationDoc struct {
			Email     string    `bson:"email"`
			OTP       string    `bson:"otp"`
			ExpiresAt time.Time `bson:"expires_at"`
			Verified  bool      `bson:"verified"`
		}

		var ev emailVerificationDoc
		err := otpCol.FindOne(
			ctx,
			bson.M{
				"email":    *req.Email,
				"otp":      *req.EmailOTP,
				"verified": false,
			},
		).Decode(&ev)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid OTP for this email"})
			return
		}

		if time.Now().After(ev.ExpiresAt) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "OTP has expired. Please request a new one."})
			return
		}

		// Mark OTP as used
		_, _ = otpCol.UpdateOne(
			ctx,
			bson.M{"email": *req.Email, "otp": *req.EmailOTP},
			bson.M{"$set": bson.M{"verified": true}},
		)

		updateFields["email"] = *req.Email
	}

	if len(updateFields) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "no changes detected"})
		return
	}

	updateFields["updated_at"] = now

	_, err := usersCol.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		logger.AddSystemLog(c,
			"profile_update_failed",
			fmt.Sprintf("user_id=%s error=%v", userID, err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	logger.AddSystemLog(c,
		"profile_update_success",
		fmt.Sprintf("user_id=%s", userID),
	)

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
