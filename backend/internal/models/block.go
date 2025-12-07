package models

import "time"

type Block struct {
	ID           string        `bson:"_id,omitempty" json:"id"`
	Index        int           `bson:"index" json:"index"`
	Timestamp    time.Time     `bson:"timestamp" json:"timestamp"`
	Transactions []Transaction `bson:"transactions" json:"transactions"`
	PreviousHash string        `bson:"previous_hash" json:"previous_hash"`
	Nonce        int64         `bson:"nonce" json:"nonce"`
	Hash         string        `bson:"hash" json:"hash"`
	MerkleRoot   string        `bson:"merkle_root,omitempty" json:"merkle_root,omitempty"`
}
