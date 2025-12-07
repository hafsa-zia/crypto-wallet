package models

import "time"

type User struct {
	ID                string    `bson:"_id,omitempty" json:"id"`
	FullName          string    `bson:"full_name" json:"full_name"`
	Email             string    `bson:"email" json:"email"`
	PasswordHash      string    `bson:"password_hash" json:"-"` // for login (optional with OTP)
	CNIC              string    `bson:"cnic" json:"cnic"`
	WalletID          string    `bson:"wallet_id" json:"wallet_id"`
	PublicKey         string    `bson:"public_key" json:"public_key"`
	EncryptedPrivKey  string    `bson:"encrypted_priv_key" json:"-"`
	Beneficiaries     []string  `bson:"beneficiaries" json:"beneficiaries"`
	ZakatDeducted     float64   `bson:"zakat_deducted" json:"zakat_deducted"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at" json:"updated_at"`
}
