// Package helpers provides a simple API client for making HTTP requests.
package helpers

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models"
)

// APIClient handles HTTP requests with support for different authentication methods.
type ApiClient struct {
	BaseURL    string
	AuthType   string // "Bearer", "Basic", or "None"
	AuthToken  string
	HTTPClient *http.Client
}

// NewApiClient creates a new API client with default settings.
// Accepts an optional timeout parameter; defaults to 60 seconds if not provided.
func NewApiClient(baseURL, authToken, authType string, timeout ...int) *ApiClient {
	defaultTimeout := 60
	if len(timeout) > 0 {
		defaultTimeout = timeout[0]
	}

	return &ApiClient{
		BaseURL:   baseURL,
		AuthToken: authToken,
		AuthType:  authType,
		HTTPClient: &http.Client{
			Timeout: time.Second * time.Duration(defaultTimeout),
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
}

// request makes an HTTP request with the specified method, endpoint, and body.
func (client *ApiClient) Request(method string, endpoint string, body interface{}) models.ApiResponse {
	// Create a new buffer for the request body
	var requestBody *bytes.Buffer
	if body != nil && (method == "POST" || method == "PUT") {
		// Marshal the body to JSON
		jsonData, err := json.Marshal(body)
		if err != nil {
			errorMessage := fmt.Sprintf("Api request got an error marshalling body to JSON: %v", err)
			return models.NewApiResponse(false, 0, nil, errorMessage, nil)
		}
		requestBody = bytes.NewBuffer(jsonData)
	} else {
		requestBody = bytes.NewBuffer(nil)
	}

	// Create a new HTTP request
	req, err := http.NewRequest(method, client.BaseURL+endpoint, requestBody)
	if err != nil {
		errorMessage := fmt.Sprintf("Api request got an error creating the request: %v", err)
		return models.NewApiResponse(false, 0, nil, errorMessage, nil)
	}

	// Set the appropriate headers for authentication
	switch client.AuthType {
	case "Bearer":
		req.Header.Add("Authorization", "Bearer "+client.AuthToken)
	case "Basic":
		encodedToken := base64.StdEncoding.EncodeToString([]byte(client.AuthToken))
		req.Header.Add("Authorization", "Basic "+encodedToken)
	}

	// Set the Content-Type header for POST and PUT requests
	if method == "POST" || method == "PUT" {
		req.Header.Set("Content-Type", "application/json")
	}

	// Send the request
	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		errorMessage := fmt.Sprintf("Api request got an error sending the request: %v", err)
		return models.NewApiResponse(false, 0, nil, errorMessage, nil)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		errorMessage := fmt.Sprintf("Api request got an error reading response body: %v", err)
		return models.NewApiResponse(false, resp.StatusCode, resp.Header, errorMessage, nil)
	}

	// Create a new ApiResponse object
	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	message := "ok"
	if !success {
		message = "ko"
	}
	return models.NewApiResponse(success, resp.StatusCode, resp.Header, message, responseBody)
}
