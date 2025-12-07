package db

import (
	"context"
	"time"

	"github.com/hafsa-zia/crypto-wallet-backend/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var DB *mongo.Database

func ConnectMongo() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		return err
	}

	Client = client
	DB = client.Database(config.AppConfig.DBName)
	return nil
}

func Col(name string) *mongo.Collection {
	return DB.Collection(name)
}
