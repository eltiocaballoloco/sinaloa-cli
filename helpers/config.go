package helpers

import (
    "fmt"
    "os"
	"io/ioutil"
	"encoding/json"
	"strings"
	"path/filepath"

    "github.com/joho/godotenv"
)

type Config struct {
    VERSION         string
    STORJ_SECRET	string
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

    // Get the version using the filePath
    version, err := GetVersion()
    if err != nil {
        return
    }

	// Set values to AppConfig
	AppConfig = Config {
		VERSION:	version,
		STORJ_SECRET: os.Getenv("STORJ_SECRET"),
	}
}

func GetVersion() (string, error) {
    // Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("\n[ERROR] Getting working directory:", err, "\n")
		return "", err
	}

	// Remove the part from cwd
	cwd = strings.TrimSuffix(cwd, "/helpers")

	// Construct the file path to version.json
	filePath := filepath.Join(cwd, "version.json")

	// Open the JSON file
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    // Read the file content
    bytes, err := ioutil.ReadAll(file)
    if err != nil {
        return "", err
    }

    // Unmarshal into a map
    var data map[string]string
    err = json.Unmarshal(bytes, &data)
    if err != nil {
        return "", err
    }

    // Return the version value
    version, ok := data["version"]
    if !ok {
        return "", fmt.Errorf("\n[ERROR] Version key not found in JSON\n")
    }

	// Return the version
    return version, nil
}
