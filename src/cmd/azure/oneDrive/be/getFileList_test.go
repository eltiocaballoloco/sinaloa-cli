package be

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/azure"

	"github.com/stretchr/testify/assert"
)

// Mocking dependencies
type MockGraphApiClient struct{}

func (m *MockGraphApiClient) GetAccessToken() (string, error) {
	if helpers.AppConfig.AZURE_CLIENT_ID == "error" {
		return "", errors.New("mock error: failed to get access token")
	}
	return "mock-access-token", nil
}

func mockLoadConfig(clientID, clientSecret, tenantID, driveID string) {
	helpers.AppConfig = helpers.Config{
		AZURE_CLIENT_ID:     clientID,
		AZURE_CLIENT_SECRET: clientSecret,
		AZURE_TENANT_ID:     tenantID,
		AZURE_DRIVE_ID:      driveID,
	}
}

func mockApiClientRequest(method, endpoint string, body interface{}) models.ApiResponse {
	if endpoint == "error" {
		return models.NewApiResponse(false, 500, nil, "Mock error: API request failed", nil)
	}

	mockResponse := azure.OneDriveGraphResponseApiModel{
		Value: []azure.OneDriveItemModel{
			{
				Name: "file1.txt",
				Type: "item",
				ParentReference: azure.ParentReference{
					Path: "/mock-path",
				},
				Folder: nil,
			},
			{
				Name: "folder1",
				Type: "folder",
				ParentReference: azure.ParentReference{
					Path: "/mock-path",
				},
				Folder: &azure.Folder{},
			},
		},
	}
	responseBody, _ := json.Marshal(mockResponse)
	return models.NewApiResponse(true, 200, nil, "Mock success", responseBody)
}

// Test cases
func TestGetDriveItems_Success(t *testing.T) {
	// Arrange: Mock dependencies
	mockLoadConfig("mock-client-id", "mock-client-secret", "mock-tenant-id", "mock-drive-id")

	// Act: Call GetDriveItems
	result, _ := GetDriveItems("/mock-path")

	// Assert: Verify response
	assert.NotEmpty(t, result, "Response should not empty")
}

func TestGetDriveItems_AccessTokenError(t *testing.T) {
	// Arrange: Mock dependencies
	mockLoadConfig("error", "mock-client-secret", "mock-tenant-id", "mock-drive-id")

	// Act: Call GetDriveItems
	result, err := GetDriveItems("/mock-path")

	// Assert: Verify response
	assert.Error(t, err, "GetDriveItems should return an error if access token retrieval fails")
	assert.Empty(t, string(result.Body), "Response should be empty")
}

func TestGetDriveItems_UnmarshalError(t *testing.T) {
	// Arrange: Mock dependencies
	mockLoadConfig("mock-client-id", "mock-client-secret", "mock-tenant-id", "mock-drive-id")

	// Act: Call GetDriveItems
	result, err := GetDriveItems("/mock-path")

	// Assert: Verify response
	assert.Error(t, err, "GetDriveItems should return an error if unmarshalling fails")
	assert.Empty(t, string(result.Body), "Response should be empty")
}

func TestGetDriveItems_MissingConfig(t *testing.T) {
	// Arrange: Mock empty configuration
	mockLoadConfig("", "", "", "")

	// Act: Call GetDriveItems
	result, err := GetDriveItems("/mock-path")

	// Assert: Verify response
	assert.Error(t, err, "GetDriveItems should return an error if configuration is missing")
	assert.Empty(t, string(result.Body), "Response should be empty")
}
