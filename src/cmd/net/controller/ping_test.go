package controller_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net/controller"

	"github.com/stretchr/testify/assert"
)

// Define a type for the ping function
type PingFunction func(string) (int, error)

// InjectedBe will allow injecting a custom Ping function for testing (simulates the behavior of the controller.Ping function)
var InjectedBe PingFunction

func init() {
	InjectedBe = func(url string) (int, error) {
		return int(200), nil
	}
}

func TestPing_Success(t *testing.T) {
	// Arrange: Mock be.Ping to simulate success
	InjectedBe = func(url string) (int, error) {
		assert.Equal(t, "google.com", url, "be.Ping should receive the correct URL")
		return http.StatusOK, nil
	}

	// Act: Call the Ping function
	response, err := controller.Ping("google.com")

	// Assert: Verify the response and no errors
	assert.NoError(t, err, "Ping should not return an error on success")
	assert.Contains(t, string(response), `200`, "Response should contain the correct status code")
	assert.Contains(t, string(response), `"Ping successful"`, "Response should indicate a successful ping")
}

func TestPing_Failure(t *testing.T) {
	// Arrange: Mock be.Ping to simulate failure
	InjectedBe = func(url string) (int, error) {
		assert.Equal(t, "thisdomaindoesntexistitislikeanexample.com", url, "be.Ping should receive the correct URL")
		return 0, errors.New("network error")
	}

	// Act: Call the Ping function
	response, err := controller.Ping("thisdomaindoesntexistitislikeanexample.com")

	// Assert: Verify the error handling
	assert.Error(t, err, "Ping should return an error on failure")
	assert.Contains(t, string(response), `no such host`, "Response should indicate a ping error")
	assert.Contains(t, string(response), `error`, "Response should contain the error message")
}

func TestPing_InvalidURL(t *testing.T) {
	// Arrange: Mock be.Ping to simulate invalid URL handling
	InjectedBe = func(url string) (int, error) {
		assert.Equal(t, "invalid-url", url, "be.Ping should receive the invalid URL")
		return 0, errors.New("invalid URL format")
	}

	// Act: Call the Ping function
	response, err := controller.Ping("invalid-url")

	// Assert: Verify the error handling
	assert.Error(t, err, "Ping should return an error for invalid URLs")
	assert.Contains(t, string(response), `"500"`, "Response should have a 500 status code")
}
