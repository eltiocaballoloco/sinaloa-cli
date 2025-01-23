package be

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/azure"
)

// UploadItem uploads a local file to a specified path in OneDrive
func UploadItem(localPath string, pathToUpload string) (bool, error) {
	// Load environment variables
	helpers.LoadConfig()

	// Initialize the Graph API client
	graphApiClient := shared.NewGraphApiClient(
		helpers.AppConfig.AZURE_CLIENT_ID,
		helpers.AppConfig.AZURE_CLIENT_SECRET,
		helpers.AppConfig.AZURE_TENANT_ID,
	)

	// Obtain an access token
	accessToken, err := graphApiClient.GetAccessToken()
	if err != nil {
		log.Fatalf("Error obtaining access token UploadItem: %v", err)
		return false, err
	}

	// Separate the directory and file name from the pathToUpload
	dir, file := filepath.Split(pathToUpload)
	if file == "" {
		err := fmt.Errorf("invalid pathToUpload: file name is missing")
		log.Fatalf("Error: %v", err)
		return false, err
	}

	// Create the upload session
	uploadSessionUrl, err := CreateUploadSession(
		graphApiClient.BaseURL+"drives/",
		accessToken,
		helpers.AppConfig.AZURE_DRIVE_ID,
		dir,
		file,
	)
	if err != nil {
		log.Fatalf("Error creating upload session (UploadOneDrive): %v", err)
		return false, err
	}

	// Upload the file in chunks
	err = UploadFileInChunks(uploadSessionUrl, localPath)
	if err != nil {
		log.Fatalf("Error uploading file in chunks (UploadOneDrive): %v", err)
		return false, err
	}

	// Print success message and return	true if all steps are successful
	fmt.Println("[INFO] File uploaded successfully!")
	return true, nil
}

// createUploadSession initiates an upload session and returns the upload URL
func CreateUploadSession(baseUrl, accessToken, driveID, folderPath, fileName string) (string, error) {
	// Declare variables
	var sessionUpload azure.OneDriveUploadSessionModel

	// Construct the URL path (endpoint) for creating an upload session
	urlPath := fmt.Sprintf("%s/root:/%s/%s:/createUploadSession", driveID, folderPath, fileName)

	// Define the request body for the post request
	requestBody := map[string]interface{}{
		"@microsoft.graph.conflictBehavior": "rename | fail | replace",
		"fileSystemInfo": map[string]interface{}{
			"@odata.type": "microsoft.graph.fileSystemInfo",
		},
		"name": fileName,
	}

	// Istance the api client to make the api call
	apiClient := helpers.NewApiClient(baseUrl, accessToken, "Bearer")

	// Use the existing request method from ApiClient to make the post request
	apiResponse := apiClient.Request(
		"POST",
		urlPath,
		requestBody,
	)

	// Parse the response
	err := json.Unmarshal(apiResponse.Body, &sessionUpload)
	if err != nil {
		return "", err
	}

	// Extract the upload URL
	if sessionUpload.UploadUrl == "" {
		return "", fmt.Errorf("Upload URL not found in response")
	}

	return sessionUpload.UploadUrl, nil
}

// uploadFileInChunks uploads the file to the provided upload URL in chunks.
// The Microsoft Graph API allows uploading
// file chunks up to a maximum of 60 MiB (megabytes) per request
func UploadFileInChunks(uploadURL, localPath string) error {
	// Open the local file
	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}
	fileSize := fileInfo.Size()

	// Determine chunk size based on file size
	var chunkSize int64
	switch {
	case fileSize <= 10*1024*1024: // Files up to 10 MiB
		chunkSize = 5 * 1024 * 1024 // 5 MiB
	case fileSize <= 100*1024*1024: // Files up to 100 MiB
		chunkSize = 10 * 1024 * 1024 // 10 MiB
	default: // Files larger than 100 MiB
		chunkSize = 20 * 1024 * 1024 // 20 MiB
	}

	// Ensure chunk size is a multiple of 320 KiB
	if chunkSize%(320*1024) != 0 {
		chunkSize = (chunkSize / (320 * 1024)) * (320 * 1024)
	}

	// Initialize variables for the upload process
	var start, end int64
	for start = 0; start < fileSize; start = end + 1 {
		end = start + chunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}

		// Read the chunk into buffer
		buffer := make([]byte, end-start+1)
		_, err := file.ReadAt(buffer, start)
		if err != nil && err != io.EOF {
			return fmt.Errorf("error reading file chunk: %v", err)
		}

		// Create a new HTTP PUT request for the chunk
		req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(buffer))
		if err != nil {
			return fmt.Errorf("error creating PUT request: %v", err)
		}
		req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))

		// Execute the request
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error uploading chunk: %v", err)
		}
		defer resp.Body.Close()

		// Check for successful response
		if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("failed to upload chunk, status: %s", resp.Status)
		}
	}

	return nil
}
