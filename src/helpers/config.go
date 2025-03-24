package helpers

import (
	"os"
)

type Config struct {
	AZURE_TENANT_ID     string
	AZURE_CLIENT_ID     string
	AZURE_CLIENT_SECRET string
	AZURE_DRIVE_ID      string
}

var (
	AppConfig Config
)

func LoadConfig() {
	// Set values to AppConfig
	AppConfig = Config{
		AZURE_TENANT_ID:     os.Getenv("AZURE_TENANT_ID"),
		AZURE_CLIENT_ID:     os.Getenv("AZURE_CLIENT_ID"),
		AZURE_CLIENT_SECRET: os.Getenv("AZURE_CLIENT_SECRET"),
		AZURE_DRIVE_ID:      os.Getenv("AZURE_DRIVE_ID"),
	}
}
