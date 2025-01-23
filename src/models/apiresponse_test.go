package models_test

import (
	"net/http"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models"

	"github.com/stretchr/testify/assert"
)

func TestNewApiResponse(t *testing.T) {
	// Arrange: Set up test data
	response := true
	statusCode := 200
	headers := http.Header{
		"Content-Type":    {"application/json"},
		"X-Custom-Header": {"custom-value"},
	}
	message := "Success"
	body := []byte(`{"key":"value"}`)

	// Act: Call the constructor
	apiResponse := models.NewApiResponse(response, statusCode, headers, message, body)

	// Assert: Verify the fields are set correctly
	assert.Equal(t, response, apiResponse.Response, "Response field should match the input")
	assert.Equal(t, statusCode, apiResponse.StatusCode, "StatusCode field should match the input")
	assert.Equal(t, headers, apiResponse.Headers, "Headers field should match the input")
	assert.Equal(t, message, apiResponse.Message, "Message field should match the input")
	assert.Equal(t, body, apiResponse.Body, "Body field should match the input")
}

func TestNewApiResponse_EmptyBody(t *testing.T) {
	// Arrange: Set up test data with an empty body
	response := false
	statusCode := 404
	headers := http.Header{
		"Content-Type": {"text/plain"},
	}
	message := "Not Found"
	var body []byte // nil body

	// Act: Call the constructor
	apiResponse := models.NewApiResponse(response, statusCode, headers, message, body)

	// Assert: Verify the fields are set correctly
	assert.Equal(t, response, apiResponse.Response, "Response field should match the input")
	assert.Equal(t, statusCode, apiResponse.StatusCode, "StatusCode field should match the input")
	assert.Equal(t, headers, apiResponse.Headers, "Headers field should match the input")
	assert.Equal(t, message, apiResponse.Message, "Message field should match the input")
	assert.Nil(t, apiResponse.Body, "Body field should be nil")
}

func TestNewApiResponse_EmptyHeaders(t *testing.T) {
	// Arrange: Set up test data with empty headers
	response := true
	statusCode := 500
	headers := http.Header{} // Empty headers
	message := "Internal Server Error"
	body := []byte("Error details")

	// Act: Call the constructor
	apiResponse := models.NewApiResponse(response, statusCode, headers, message, body)

	// Assert: Verify the fields are set correctly
	assert.Equal(t, response, apiResponse.Response, "Response field should match the input")
	assert.Equal(t, statusCode, apiResponse.StatusCode, "StatusCode field should match the input")
	assert.Equal(t, headers, apiResponse.Headers, "Headers field should match the input")
	assert.Equal(t, message, apiResponse.Message, "Message field should match the input")
	assert.Equal(t, body, apiResponse.Body, "Body field should match the input")
}
