package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmailOTP struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email     string             `bson:"email" json:"email"`
	OTP       string             `bson:"otp" json:"-"` // never expose OTP in JSON
	ExpiresAt time.Time          `bson:"expires_at" json:"expires_at"`
	Verified  bool               `bson:"verified" json:"verified"`
}
