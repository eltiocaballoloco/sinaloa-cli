package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	VERSION      string
	STORJ_SECRET string
}

var (
	AppConfig Config
)

func LoadConfig() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		return
	}

	// Set values to AppConfig
	AppConfig = Config{
		STORJ_SECRET: os.Getenv("STORJ_SECRET"),
	}
}
