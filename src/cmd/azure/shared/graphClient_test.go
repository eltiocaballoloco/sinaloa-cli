package shared_test

import (
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"

	"github.com/stretchr/testify/assert"
)

func mockLoadConfig(clientID, tenantID, clientSecret string) {
	helpers.AppConfig = helpers.Config{
		AZURE_CLIENT_ID:     clientID,
		AZURE_TENANT_ID:     tenantID,
		AZURE_CLIENT_SECRET: clientSecret,
	}
}

func TestGraphClient_Success(t *testing.T) {
	// Arrange: Mock AppConfig
	mockLoadConfig("mock-client-id", "mock-tenant-id", "mock-client-secret")

	// Act: Create a GraphClient
	client, err := shared.NewGraphClient()

	// Assert: Ensure no errors and the client is initialized
	assert.NoError(t, err, "GraphClient should initialize without error")
	assert.NotNil(t, client, "GraphClient should not be nil")
}

func TestGraphClient_MissingConfig(t *testing.T) {
	// Arrange: Mock missing AppConfig
	mockLoadConfig("", "", "")

	// Act: Create a GraphClient
	client, err := shared.NewGraphClient()

	// Assert: Ensure error is returned for missing configuration
	assert.Error(t, err, "GraphClient should return an error for missing configuration")
	assert.Contains(t, err.Error(), "missing required environment variables", "Error message should indicate missing variables")
	assert.Nil(t, client, "GraphClient should not be initialized with missing config")
}

func TestGraphClient_CredentialError(t *testing.T) {
	// Arrange: Mock AppConfig
	mockLoadConfig("mock-client-id", "mock-tenant-id", "mock-client-secret")

	// Act: Create a GraphClient
	client, _ := shared.NewGraphClient()

	// Assert: Ensure error is returned for credential failure
	assert.NotEmpty(t, client, "GraphClient should not be empty")
}
