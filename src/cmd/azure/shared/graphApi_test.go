package shared_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/shared"

	"github.com/stretchr/testify/assert"
)

func TestNewGraphApiClient(t *testing.T) {
	// Act: Create a new GraphApiClient
	client := shared.NewGraphApiClient("mock-client-id", "mock-client-secret", "mock-tenant-id")

	// Assert: Verify fields are set correctly
	assert.Equal(t, "mock-client-id", client.ClientID)
	assert.Equal(t, "mock-client-secret", client.ClientSecret)
	assert.Equal(t, "mock-tenant-id", client.TenantID)
	assert.Equal(t, "https://graph.microsoft.com/v1.0/", client.BaseURL)
	assert.NotNil(t, client.HTTPClient, "HTTPClient should be initialized")
}

func TestGraphApiClient_GetAccessToken_Success(t *testing.T) {
	// Arrange: Set up a mock HTTP server to simulate a successful response
	mockResponse := `{"access_token": "mock-access-token"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method, "Request method should be POST")
		assert.Contains(t, r.Header.Get("Content-Type"), "application/x-www-form-urlencoded", "Content-Type header should be set")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Replace BaseURL with the mock server URL
	client := &shared.GraphApiClient{
		BaseURL:      server.URL,
		HTTPClient:   &http.Client{},
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		TenantID:     "mock-tenant-id",
	}

	// Act: Call GetAccessToken
	token, _ := client.GetAccessToken()

	// Assert: Verify the access token and no error
	assert.NotEmpty(t, "mock-access-token", token, "Access token should match the mock response")
}

func TestGraphApiClient_GetAccessToken_InvalidResponse(t *testing.T) {
	// Arrange: Set up a mock HTTP server to simulate an invalid response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"invalid_key": "invalid-value"}`)) // No access_token field
	}))
	defer server.Close()

	// Replace BaseURL with the mock server URL
	client := &shared.GraphApiClient{
		BaseURL:      server.URL,
		HTTPClient:   &http.Client{},
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		TenantID:     "mock-tenant-id",
	}

	// Act: Call GetAccessToken
	token, err := client.GetAccessToken()

	// Assert: Verify error and empty token
	assert.Error(t, err, "GetAccessToken should return an error if access_token is missing")
	assert.Empty(t, token, "Access token should be empty on invalid response")
	assert.Contains(t, err.Error(), "token response did not contain access_token", "Error message should indicate missing access_token")
}

func TestGraphApiClient_GetAccessToken_HTTPError(t *testing.T) {
	// Arrange: Set up a mock HTTP server to simulate a failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "mock server error", http.StatusInternalServerError)
	}))
	defer server.Close()

	// Replace BaseURL with the mock server URL
	client := &shared.GraphApiClient{
		BaseURL:      server.URL,
		HTTPClient:   &http.Client{},
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		TenantID:     "mock-tenant-id",
	}

	// Act: Call GetAccessToken
	token, err := client.GetAccessToken()

	// Assert: Verify error and empty token
	assert.Error(t, err, "GetAccessToken should return an error on HTTP error")
	assert.Empty(t, token, "Access token should be empty on HTTP error")
}

func TestGraphApiClient_GetAccessToken_DecodeError(t *testing.T) {
	// Arrange: Set up a mock HTTP server to simulate a decoding error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`invalid-json`)) // Invalid JSON
	}))
	defer server.Close()

	// Replace BaseURL with the mock server URL
	client := &shared.GraphApiClient{
		BaseURL:      server.URL,
		HTTPClient:   &http.Client{},
		ClientID:     "mock-client-id",
		ClientSecret: "mock-client-secret",
		TenantID:     "mock-tenant-id",
	}

	// Act: Call GetAccessToken
	token, err := client.GetAccessToken()

	// Assert: Verify error and empty token
	assert.Error(t, err, "GetAccessToken should return an error on JSON decode failure")
	assert.Empty(t, token, "Access token should be empty on decode failure")
}
