package utxo

import (
	"context"
	"errors"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetBalance(ctx context.Context, walletID string) (float64, error) {
	col := db.Col("utxos")
	cur, err := col.Find(ctx, bson.M{"owner_wallet": walletID, "is_spent": false})
	if err != nil {
		return 0, err
	}
	defer cur.Close(ctx)

	var u models.UTXO
	var total float64
	for cur.Next(ctx) {
		if err := cur.Decode(&u); err != nil {
			return 0, err
		}
		total += u.Amount
	}
	return total, nil
}

func SelectUTXOsForAmount(ctx context.Context, walletID string, amount float64) ([]models.UTXO, float64, error) {
	col := db.Col("utxos")
	cur, err := col.Find(ctx, bson.M{"owner_wallet": walletID, "is_spent": false})
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx)

	var chosen []models.UTXO
	var sum float64
	for cur.Next(ctx) {
		var u models.UTXO
		if err := cur.Decode(&u); err != nil {
			return nil, 0, err
		}
		chosen = append(chosen, u)
		sum += u.Amount
		if sum >= amount {
			change := sum - amount
			return chosen, change, nil
		}
	}
	return nil, 0, errors.New("insufficient funds")
}
