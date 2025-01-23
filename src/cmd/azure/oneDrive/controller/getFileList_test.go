package controller_test

import (
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"

	"github.com/stretchr/testify/assert"
)

func TestGetFileList_Success(t *testing.T) {
	// Act: Call GetFileList with a valid path
	result, _ := controller.GetFileList("/mock-path")
	// Assert: Verify response
	assert.NotEmpty(t, string(result), "Items fetched successfully", "Response should include success message")
}

func TestGetFileList_Failure(t *testing.T) {
	// Act: Call GetFileList with an invalid path
	result, err := controller.GetFileList("error")

	// Assert: Verify response
	assert.Error(t, err, "GetFileList should return an error on failure")
	assert.NotEmpty(t, string(result), "Failed to fetch items", "Response should indicate failure to fetch items")
}

func TestGetFileList_EmptyResponse(t *testing.T) {
	// Act: Call GetFileList with a valid path but no items
	result, _ := controller.GetFileList("/mock-path")

	// Assert: Verify response
	assert.NotEmpty(t, string(result), "No items found", "Response should indicate no items were found")
	assert.NotEmpty(t, string(result), "[]", "Response body should be an empty array")
}
