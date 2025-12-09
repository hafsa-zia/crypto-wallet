package zakat

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hafsa-zia/crypto-wallet-backend/internal/db"
	"github.com/hafsa-zia/crypto-wallet-backend/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const ZakatWalletID = "ZAKAT_POOL"

func RunZakatForWallet(c *gin.Context, walletID string) (bson.M, error) {
	ctx := context.Background()

	utxosCol := db.Col("utxos")
	txsCol := db.Col("transactions")
	usersCol := db.Col("users")

	// TODO: adjust "owner_wallet" + "spent" to match YOUR utxo schema
	cur, err := utxosCol.Find(ctx, bson.M{
		"owner_wallet": walletID,
		"spent":        false,
	})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var utxos []bson.M
	if err := cur.All(ctx, &utxos); err != nil {
		return nil, err
	}
	if len(utxos) == 0 {
		return nil, errors.New("zakat not due")
	}

	var balance float64
	for _, u := range utxos {
		if amt, ok := u["amount"].(float64); ok {
			balance += amt
		}
	}
	if balance <= 0 {
		return nil, errors.New("balance <= 0, zakat not due")
	}

	zakatAmount := balance * 0.025
	if zakatAmount <= 0 {
		return nil, errors.New("zakat computed as 0, nothing to do")
	}

	var inputs []bson.M
	var inputSum float64

	for _, u := range utxos {
		rawID, hasID := u["_id"]
		if !hasID {
			continue
		}
		oid, ok := rawID.(primitive.ObjectID)
		if !ok {
			continue
		}
		amt, _ := u["amount"].(float64)

		inputs = append(inputs, bson.M{
			"utxo_id": oid.Hex(),
		})
		inputSum += amt

		if inputSum >= zakatAmount {
			break
		}
	}

	if inputSum < zakatAmount {
		return nil, errors.New("not enough UTXOs to cover zakat")
	}

	change := inputSum - zakatAmount

	var outputs []bson.M
	outputs = append(outputs, bson.M{
		"owner_wallet": ZakatWalletID,
		"amount":       zakatAmount,
	})
	if change > 0 {
		outputs = append(outputs, bson.M{
			"owner_wallet": walletID,
			"amount":       change,
		})
	}

	now := time.Now().UTC()

	txDoc := bson.M{
		"sender_wallet":   walletID,
		"receiver_wallet": ZakatWalletID,
		"amount":          zakatAmount,
		"note":            "Monthly zakat deduction (2.5%)",
		"type":            "zakat_deduction",
		"status":          "pending",
		"timestamp":       now,
		"inputs":          inputs,
		"outputs":         outputs,
	}

	res, err := txsCol.InsertOne(ctx, txDoc)
	if err != nil {
		return nil, err
	}

	txDoc["_id"] = res.InsertedID

	_, _ = usersCol.UpdateOne(ctx, bson.M{"wallet_id": walletID}, bson.M{
		"$inc": bson.M{"zakat_deducted": zakatAmount},
	})

	logger.AddSystemLog(c,
		"zakat_created",
		fmt.Sprintf("wallet=%s amount=%.4f", walletID, zakatAmount),
	)

	return txDoc, nil
}
