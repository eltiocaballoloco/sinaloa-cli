package response_test

import (
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/response"

	"github.com/stretchr/testify/assert"
)

func TestNewResponse(t *testing.T) {
	// Arrange: Set up test data
	expectedResponse := true
	expectedCode := "200"
	expectedMessage := "Success"
	expectedData := map[string]interface{}{
		"id":   1,
		"name": "Test Name",
	}

	// Act: Call the function under test
	resp := response.NewResponse(expectedResponse, expectedCode, expectedMessage, expectedData)

	// Assert: Validate the returned struct
	assert.Equal(t, expectedResponse, resp.Response, "Response field should match")
	assert.Equal(t, expectedCode, resp.Code, "Code field should match")
	assert.Equal(t, expectedMessage, resp.Message, "Message field should match")
	assert.Equal(t, expectedData, resp.Data, "Data field should match")
}

func TestResponseEmptyData(t *testing.T) {
	// Arrange: No data for the response
	expectedResponse := false
	expectedCode := "404"
	expectedMessage := "Not Found"
	var expectedData interface{} = nil

	// Act: Call the function under test
	resp := response.NewResponse(expectedResponse, expectedCode, expectedMessage, expectedData)

	// Assert: Validate the returned struct
	assert.Equal(t, expectedResponse, resp.Response, "Response field should match")
	assert.Equal(t, expectedCode, resp.Code, "Code field should match")
	assert.Equal(t, expectedMessage, resp.Message, "Message field should match")
	assert.Nil(t, resp.Data, "Data field should be nil")
}
