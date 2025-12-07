package models

import "time"

// ZakatRecord stores monthly zakat deductions per wallet (for reporting)
type ZakatRecord struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	WalletID  string    `bson:"wallet_id" json:"wallet_id"`
	Amount    float64   `bson:"amount" json:"amount"`
	Month     string    `bson:"month" json:"month"` // e.g. "2025-12"
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}
