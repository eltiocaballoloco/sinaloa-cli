package controller

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

func GetFile(path string, pathToSaveFile string) ([]byte, error) {
	// Extract the directory from the path
	directoryPath := filepath.Dir(path)

	// Split the path to get the required file name
	requiredName := filepath.Base(path)

	// Call the GetDriveItems function from the backend
	apiResponse, err := be.GetDriveItems(directoryPath)
	if err != nil {
		return helpers.HandleController(
			false,
			"500",
			"Failed to fetch items",
			"GetFile",
			nil,
			fmt.Errorf("failed to fetch items: %v", err),
		)
	}

	// Parse the raw JSON data into a generic map
	var rawData map[string]interface{}
	err = json.Unmarshal(apiResponse.Body, &rawData)
	if err != nil {
		return helpers.HandleController(
			false,
			"500",
			"Failed to parse JSON",
			"GetFile",
			nil,
			fmt.Errorf("failed to parse JSON: %v", err),
		)
	}

	// Extract the "values" key from the parsed JSON
	values, ok := rawData["values"].([]interface{})
	if !ok {
		return helpers.HandleController(
			false,
			"404",
			"No items found in the response",
			"GetFile",
			nil,
			fmt.Errorf("no items found in the response"),
		)
	}

	// Iterate over the values
	for _, v := range values {
		// Check if the item is a map
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if the item has the required name and is of type "item"
		if item["name"] == requiredName && item["type"] == "item" {
			// If "no_store", return the file metadata
			if pathToSaveFile == "no_store" {
				return helpers.HandleController(
					apiResponse.Response,
					fmt.Sprintf("%d", apiResponse.StatusCode),
					"File metadata retrieved successfully",
					"GetFile",
					item,
					nil,
				)
			}

			// Otherwise, download the file
			downloadUrl, ok := item["downloadUrl"].(string)
			if !ok || downloadUrl == "" {
				return helpers.HandleController(
					false,
					"404",
					fmt.Sprintf("File '%s' does not have a download URL", requiredName),
					"GetFile",
					nil,
					fmt.Errorf("file '%s' does not have a download URL", requiredName),
				)
			}

			err := downloadFile(downloadUrl, pathToSaveFile)
			if err != nil {
				return helpers.HandleController(
					false,
					"500",
					"Failed to download file",
					"GetFile",
					nil,
					fmt.Errorf("failed to download file: %v", err),
				)
			}

			// Return success message
			return helpers.HandleController(
				apiResponse.Response,
				fmt.Sprintf("%d", apiResponse.StatusCode),
				"File downloaded successfully",
				"GetFile",
				item,
				nil,
			)
		}
	}

	// If no match is found, return a structured error response
	return helpers.HandleController(
		false,
		"404",
		fmt.Sprintf("File with name '%s' not found in path '%s'", requiredName, path),
		"GetFile",
		nil,
		fmt.Errorf("file with name '%s' not found in path '%s'", requiredName, path),
	)
}

func downloadFile(downloadUrl, savePath string) error {
	// Create the HTTP request
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: HTTP status %d", resp.StatusCode)
	}

	// Create the file on disk
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the file content
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file content: %v", err)
	}

	return nil
}
