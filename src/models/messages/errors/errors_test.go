package errors_test

import (
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewErrorResponse(t *testing.T) {
	// Arrange: Set up expected values
	expectedResponse := false
	expectedCode := "400"
	expectedMessage := "Bad Request"

	// Act: Call the function under test
	errResp := errors.NewErrorResponse(expectedResponse, expectedCode, expectedMessage)

	// Assert: Verify the struct fields
	assert.Equal(t, expectedResponse, errResp.Response, "Response field should match")
	assert.Equal(t, expectedCode, errResp.Code, "Code field should match")
	assert.Equal(t, expectedMessage, errResp.Message, "Message field should match")
	assert.Equal(t, struct{}{}, errResp.Data, "Data field should be an empty struct")
}

func TestNewErrorResponseWithEmptyMessage(t *testing.T) {
	// Arrange: Set up values with an empty custom message
	expectedResponse := false
	expectedCode := "500"
	expectedMessage := ""

	// Act: Call the function under test
	errResp := errors.NewErrorResponse(expectedResponse, expectedCode, expectedMessage)

	// Assert: Verify the struct fields
	assert.Equal(t, expectedResponse, errResp.Response, "Response field should match")
	assert.Equal(t, expectedCode, errResp.Code, "Code field should match")
	assert.Equal(t, expectedMessage, errResp.Message, "Message field should match")
	assert.Equal(t, struct{}{}, errResp.Data, "Data field should be an empty struct")
}
