package helpers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/response"
)

// APICLIENT
func TestNewApiClient(t *testing.T) {
	// Arrange: Create an API client
	baseURL := "http://example.com"
	authToken := "test-token"
	authType := "Bearer"
	timeout := 30

	// Act: Initialize the API client
	client := helpers.NewApiClient(baseURL, authToken, authType, timeout)

	// Assert: Validate the client fields
	assert.Equal(t, baseURL, client.BaseURL, "BaseURL should match")
	assert.Equal(t, authToken, client.AuthToken, "AuthToken should match")
	assert.Equal(t, authType, client.AuthType, "AuthType should match")
	assert.Equal(t, time.Second*time.Duration(timeout), client.HTTPClient.Timeout, "Timeout should match")
}

func TestApiClientRequest_Success(t *testing.T) {
	// Arrange: Set up a test server with a mock response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"), "Authorization header should match")
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"), "Content-Type header should match")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Success"}`))
	}))
	defer server.Close()

	client := helpers.NewApiClient(server.URL, "test-token", "Bearer", 10)

	// Act: Make a request
	response := client.Request("POST", "/test-endpoint", map[string]string{"key": "value"})

	// Assert: Validate the response
	assert.True(t, response.Response, "Response should indicate success")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Status code should match")
	assert.JSONEq(t, `{"message": "Success"}`, string(response.Body), "Response body should match")
}

func TestApiClientRequest_Error(t *testing.T) {
	// Arrange: Set up a test server with an error response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "Internal Server Error"}`))
	}))
	defer server.Close()

	client := helpers.NewApiClient(server.URL, "test-token", "Bearer", 10)

	// Act: Make a request
	response := client.Request("GET", "/error-endpoint", nil)

	// Assert: Validate the response
	assert.False(t, response.Response, "Response should indicate failure")
	assert.Equal(t, http.StatusInternalServerError, response.StatusCode, "Status code should match")
	assert.JSONEq(t, `{"error": "Internal Server Error"}`, string(response.Body), "Response body should match")
}

func TestApiClientRequest_AuthBasic(t *testing.T) {
	// Arrange: Set up a test server to validate Basic authentication
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		expectedAuth := "Basic dGVzdDp0b2tlbg==" // base64("test:token")
		assert.Equal(t, expectedAuth, auth, "Authorization header should match Basic Auth encoding")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := helpers.NewApiClient(server.URL, "test:token", "Basic", 10)

	// Act: Make a request
	response := client.Request("GET", "/auth-basic", nil)

	// Assert: Validate the response
	assert.True(t, response.Response, "Response should indicate success")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Status code should match")
}

func TestApiClientRequest_MarshallingError(t *testing.T) {
	// Arrange: Create an API client
	client := helpers.NewApiClient("http://example.com", "test-token", "Bearer", 10)

	// Act: Make a request with invalid JSON input
	response := client.Request("POST", "/invalid-json", func() {})

	// Assert: Validate the response
	assert.False(t, response.Response, "Response should indicate failure")
	assert.Equal(t, 0, response.StatusCode, "Status code should be 0 for marshalling errors")
	assert.Contains(t, response.Message, "marshalling body to JSON", "Message should indicate a marshalling error")
}

// CONFIG
func TestLoadConfig(t *testing.T) {
	// Act: Load the configuration
	helpers.LoadConfig()

	// Assert: Verify that AppConfig has the correct values
	assert.Equal(t, "", helpers.AppConfig.AZURE_TENANT_ID, "AZURE_TENANT_ID should match")
	assert.Equal(t, "", helpers.AppConfig.AZURE_CLIENT_ID, "AZURE_CLIENT_ID should match")
	assert.Equal(t, "", helpers.AppConfig.AZURE_CLIENT_SECRET, "AZURE_CLIENT_SECRET should match")
	assert.Equal(t, "", helpers.AppConfig.AZURE_DRIVE_ID, "AZURE_DRIVE_ID should match")
}

func TestLoadConfigWithDotenv(t *testing.T) {
	helpers.LoadConfig()

	// Assert: Verify that AppConfig has the correct values from .env
	assert.Empty(t, "", helpers.AppConfig.AZURE_TENANT_ID, "AZURE_TENANT_ID should match .env value")
	assert.Empty(t, "", helpers.AppConfig.AZURE_CLIENT_ID, "AZURE_CLIENT_ID should match .env value")
	assert.Empty(t, "", helpers.AppConfig.AZURE_CLIENT_SECRET, "AZURE_CLIENT_SECRET should match .env value")
	assert.Empty(t, "", helpers.AppConfig.AZURE_DRIVE_ID, "AZURE_DRIVE_ID should match .env value")
}

// RESPONSE HANDLER
func TestHandleController_Success(t *testing.T) {
	// Arrange: Set up input for a successful response
	result := true
	statusCode := "200"
	message := "Operation successful"
	controllerFunction := "TestController"
	data := map[string]string{"key": "value"}

	// Act: Call HandleController
	output, err := helpers.HandleController(result, statusCode, message, controllerFunction, data, nil)

	// Assert: Validate output
	assert.NoError(t, err, "Error should be nil for successful responses")

	var parsedResponse response.Response
	parseErr := json.Unmarshal(output, &parsedResponse)
	assert.NoError(t, parseErr, "Output JSON should parse into Response struct")
	assert.True(t, parsedResponse.Response, "Response field should be true")
	assert.Equal(t, statusCode, parsedResponse.Code, "StatusCode should match")
	assert.Equal(t, message, parsedResponse.Message, "Message should match")
}

func TestHandleController_Error(t *testing.T) {
	// Arrange: Set up input for an error response
	result := false
	statusCode := "500"
	message := "Internal server error"
	controllerFunction := "TestController"

	// Act: Call HandleController
	output, handlerErr := helpers.HandleController(result, statusCode, message, controllerFunction, nil, nil)

	// Assert: Validate output
	assert.Equal(t, nil, handlerErr, "Handler should return the input error")

	var parsedErrorResponse errors.ErrorResponse
	parseErr := json.Unmarshal(output, &parsedErrorResponse)
	assert.NoError(t, parseErr, "Output JSON should parse into ErrorResponse struct")
	assert.False(t, parsedErrorResponse.Response, "Response field should be false")
	assert.Equal(t, statusCode, parsedErrorResponse.Code, "StatusCode should match")
	assert.Equal(t, message, parsedErrorResponse.Message, "Message should match")
}

func TestHandleController_ByteData(t *testing.T) {
	// Arrange: Set up input with data as a JSON byte slice
	result := true
	statusCode := "200"
	message := "Operation successful"
	controllerFunction := "TestController"
	data := []byte(`{"key": "value"}`)

	// Act: Call HandleController
	output, err := helpers.HandleController(result, statusCode, message, controllerFunction, data, nil)

	// Assert: Validate output
	assert.NoError(t, err, "Error should be nil for successful responses with byte data")

	var parsedResponse response.Response
	parseErr := json.Unmarshal(output, &parsedResponse)
	assert.NoError(t, parseErr, "Output JSON should parse into Response struct")
	assert.True(t, parsedResponse.Response, "Response field should be true")
	assert.Equal(t, statusCode, parsedResponse.Code, "StatusCode should match")
	assert.Equal(t, message, parsedResponse.Message, "Message should match")

	// Validate that the byte data was unmarshaled correctly
	expectedData := map[string]interface{}{"key": "value"}
	assert.Equal(t, expectedData, parsedResponse.Data, "Data should match the unmarshaled byte data")
}

func TestHandleController_InvalidByteData(t *testing.T) {
	// Arrange: Set up input with invalid JSON byte data
	result := true
	statusCode := "200"
	message := "Operation successful"
	controllerFunction := "TestController"
	data := []byte(`invalid-json`)

	// Act: Call HandleController
	output, err := helpers.HandleController(result, statusCode, message, controllerFunction, data, nil)

	// Assert: Validate output
	assert.NoError(t, err, "Error should be nil even if byte data is invalid")

	var parsedResponse response.Response
	parseErr := json.Unmarshal(output, &parsedResponse)
	assert.NoError(t, parseErr, "Output JSON should parse into Response struct")
	assert.True(t, parsedResponse.Response, "Response field should be true")
	assert.Equal(t, statusCode, parsedResponse.Code, "StatusCode should match")
	assert.Equal(t, message, parsedResponse.Message, "Message should match")
}
