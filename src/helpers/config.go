package helpers

import (
	"os"
	"strconv"
)

type Config struct {
	SINALOA_DEBUG       bool
	AZURE_TENANT_ID     string
	AZURE_CLIENT_ID     string
	AZURE_CLIENT_SECRET string
	AZURE_DRIVE_ID      string
	DOCKER_HUB_USER_R   string
	DOCKER_HUB_PWD_R    string
}

var (
	AppConfig Config
)

func LoadConfig() {
	// Set debug mode
	debug, err := strconv.ParseBool(os.Getenv("SINALOA_DEBUG"))
	if err != nil {
		debug = false // default to false if parsing fails or not set
	}
	// Set values to AppConfig
	AppConfig = Config{
		SINALOA_DEBUG:       debug,
		AZURE_TENANT_ID:     os.Getenv("AZURE_TENANT_ID"),
		AZURE_CLIENT_ID:     os.Getenv("AZURE_CLIENT_ID"),
		AZURE_CLIENT_SECRET: os.Getenv("AZURE_CLIENT_SECRET"),
		AZURE_DRIVE_ID:      os.Getenv("AZURE_DRIVE_ID"),
		DOCKER_HUB_USER_R:   os.Getenv("DOCKER_HUB_USER_R"),
		DOCKER_HUB_PWD_R:    os.Getenv("DOCKER_HUB_PWD_R"),
	}
}
