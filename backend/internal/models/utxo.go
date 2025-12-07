package models

type UTXO struct {
	ID           string  `bson:"_id,omitempty" json:"id"`
	TxID         string  `bson:"tx_id" json:"tx_id"`
	Index        int     `bson:"index" json:"index"`
	OwnerWallet  string  `bson:"owner_wallet" json:"owner_wallet"`
	Amount       float64 `bson:"amount" json:"amount"`
	IsSpent      bool    `bson:"is_spent" json:"is_spent"`
	SpentInTxID  string  `bson:"spent_in_tx_id,omitempty" json:"spent_in_tx_id,omitempty"`
}
