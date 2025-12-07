package models

type Wallet struct {
	ID       string  `bson:"_id,omitempty" json:"id"`
	WalletID string  `bson:"wallet_id" json:"wallet_id"`
	UserID   string  `bson:"user_id" json:"user_id"`
	Balance  float64 `bson:"balance" json:"balance"` // cached, must be validated via UTXO
}
