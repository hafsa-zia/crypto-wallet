package models

import "time"

type TxUTXOInput struct {
	UTXOId string `bson:"utxo_id" json:"utxo_id"`
	Index  int    `bson:"index" json:"index"`
}

type TxUTXOOutput struct {
	OwnerWallet string  `bson:"owner_wallet" json:"owner_wallet"`
	Amount      float64 `bson:"amount" json:"amount"`
}

type Transaction struct {
	ID             string         `bson:"_id,omitempty" json:"id"`
	SenderWallet   string         `bson:"sender_wallet" json:"sender_wallet"`
	ReceiverWallet string         `bson:"receiver_wallet" json:"receiver_wallet"`
	Amount         float64        `bson:"amount" json:"amount"`
	Note           string         `bson:"note" json:"note"`
	Timestamp      time.Time      `bson:"timestamp" json:"timestamp"`
	SenderPubKey   string         `bson:"sender_public_key" json:"sender_public_key"`
	Signature      string         `bson:"signature" json:"signature"`
	Inputs         []TxUTXOInput  `bson:"inputs" json:"inputs"`
	Outputs        []TxUTXOOutput `bson:"outputs" json:"outputs"`
	Type           string         `bson:"type" json:"type"` // normal, zakat_deduction, mining_reward
	BlockID        string         `bson:"block_id,omitempty" json:"block_id,omitempty"`
	Status         string         `bson:"status" json:"status"` // pending, confirmed, rejected
}
