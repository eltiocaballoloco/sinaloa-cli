package controller_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"

	"github.com/stretchr/testify/assert"
)

func TestGetFile_Success(t *testing.T) {
	// Act: Call GetFile with valid inputs
	result, _ := controller.GetFile("file1.txt", "/tmp/file1.txt")
	// Assert: Verify response
	assert.NotEmpty(t, result, "GetFile should not be empty on success")
}

func TestGetFile_MetadataOnly(t *testing.T) {
	// Act: Call GetFile with "no_store" to get metadata only
	result, _ := controller.GetFile("file1.txt", "no_store")
	// Assert: Verify response
	assert.NotEmpty(t, result, "GetFile should not return an error when retrieving metadata")
}

func TestGetFile_FileNotFound(t *testing.T) {
	// Act: Call GetFile with a file that doesn't exist
	result, err := controller.GetFile("nonexistent.txt", "no_store")

	// Assert: Verify response
	assert.Error(t, err, "GetFile should return an error when file is not found")
	assert.NotEmpty(t, string(result), "File with name 'nonexistent.txt' not found", "Response should indicate file not found")
}

func TestGetFile_DriveItemsError(t *testing.T) {
	// Act: Call GetFile with an invalid path
	result, err := controller.GetFile("file1.txt", "/tmp/file1.txt")

	// Assert: Verify response
	assert.Error(t, err, "GetFile should return an error when GetDriveItems fails")
	assert.NotEmpty(t, string(result), "Failed to fetch items", "Response should indicate failure to fetch items")
}

func TestDownloadFile_Success(t *testing.T) {
	// Create a temporary file path
	tmpFile := filepath.Join(os.TempDir(), "mock-file.txt")
	defer os.Remove(tmpFile) // Clean up after test

	// Act: Call downloadFile
	err := controller.DownloadFile("http://mock-url.com/file1.txt", tmpFile)

	// Assert: Verify no error
	assert.NotEmpty(t, err, "downloadFile should not return an error on successful download")
}

func TestDownloadFile_Failure(t *testing.T) {
	// Act: Call downloadFile with an invalid URL
	err := controller.DownloadFile("http://mock-url.com/nonexistent.txt", "/tmp/mock-file.txt")

	// Assert: Verify error
	assert.Error(t, err, "downloadFile should return an error on failure")
	assert.NotEmpty(t, err.Error(), "mock error: failed to download file", "Error message should indicate download failure")
}
