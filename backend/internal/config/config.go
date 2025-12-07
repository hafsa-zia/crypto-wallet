package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI      string
	DBName        string
	JWTSecret     string
	AESSecretKey  string
	ZakatWalletID string
	PowDifficulty int
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load()

	diff, err := strconv.Atoi(os.Getenv("POW_DIFFICULTY"))
	if err != nil {
		diff = 5
	}

	AppConfig = &Config{
		MongoURI:      os.Getenv("MONGODB_URI"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		AESSecretKey:  os.Getenv("AES_SECRET_KEY"),
		ZakatWalletID: os.Getenv("ZAKAT_WALLET_ID"),
		PowDifficulty: diff,
	}

	if AppConfig.MongoURI == "" {
		log.Fatal("MONGODB_URI not set")
	}
}
