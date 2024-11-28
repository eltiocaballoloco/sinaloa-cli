package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AZURE_TENANT_ID     string
	AZURE_CLIENT_ID     string
	AZURE_CLIENT_SECRET string
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
		AZURE_TENANT_ID:     os.Getenv("AZURE_TENANT_ID"),
		AZURE_CLIENT_ID:     os.Getenv("AZURE_CLIENT_ID"),
		AZURE_CLIENT_SECRET: os.Getenv("AZURE_CLIENT_SECRET"),
	}
}
