package scheduler

import (
	"context"
	"time"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/models"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/utxo"
	"go.mongodb.org/mongo-driver/bson"
)

func RunZakatNow(ctx context.Context) error {
	walletsCol := db.Col("wallets")
	cur, err := walletsCol.Find(ctx, bson.M{})
	if err != nil {
		return err
	}
	defer cur.Close(ctx)

	pendingCol := db.Col("pending_transactions")
	for cur.Next(ctx) {
		var w models.Wallet
		if err := cur.Decode(&w); err != nil {
			continue
		}
		balance, err := utxo.GetBalance(ctx, w.WalletID)
		if err != nil || balance <= 0 {
			continue
		}
		zakat := balance * 0.025
		if zakat <= 0 {
			continue
		}

		tx := models.Transaction{
			SenderWallet:   w.WalletID,
			ReceiverWallet: config.AppConfig.ZakatWalletID,
			Amount:         zakat,
			Note:           "Monthly Zakat",
			Timestamp:      time.Now().UTC(),
			Type:           "zakat_deduction",
			Status:         "pending",
		}
		_, _ = pendingCol.InsertOne(ctx, tx)

		// Update user zakat tracking
		_, _ = db.Col("users").UpdateOne(ctx,
			bson.M{"wallet_id": w.WalletID},
			bson.M{"$inc": bson.M{"zakat_deducted": zakat}},
		)
	}
	return nil
}
