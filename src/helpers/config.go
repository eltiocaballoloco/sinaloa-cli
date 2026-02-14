package helpers

import (
	"os"
	"strconv"
)

type Config struct {
	SINALOA_DEBUG       bool
	ARGOCD_URL          string
	ARGOCD_USER         string
	ARGOCD_PASSWORD     string
	AZURE_TENANT_ID     string
	AZURE_CLIENT_ID     string
	AZURE_CLIENT_SECRET string
	AZURE_DRIVE_ID      string
	DOCKER_HUB_USER_RWD string
	DOCKER_HUB_PWD_RWD  string
	GITHUB_TOKEN        string
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
		ARGOCD_URL:          os.Getenv("ARGOCD_URL"),
		ARGOCD_USER:         os.Getenv("ARGOCD_USER"),
		ARGOCD_PASSWORD:     os.Getenv("ARGOCD_PASSWORD"),
		AZURE_TENANT_ID:     os.Getenv("AZURE_TENANT_ID"),
		AZURE_CLIENT_ID:     os.Getenv("AZURE_CLIENT_ID"),
		AZURE_CLIENT_SECRET: os.Getenv("AZURE_CLIENT_SECRET"),
		AZURE_DRIVE_ID:      os.Getenv("AZURE_DRIVE_ID"),
		DOCKER_HUB_USER_RWD: os.Getenv("DOCKER_HUB_USER_RWD"),
		DOCKER_HUB_PWD_RWD:  os.Getenv("DOCKER_HUB_PWD_RWD"),
		GITHUB_TOKEN:        os.Getenv("GITHUB_TOKEN"),
	}
}
