package be_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net/be"

	"github.com/stretchr/testify/assert"
)

func TestPing_SuccessWithDomain(t *testing.T) {
	// Arrange: Set up a mock server to simulate a successful ping
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "HEAD", r.Method, "Expected HEAD method")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Act: Call Ping with the mock server's URL
	statusCode, err := be.Ping(server.URL[len("http://"):])

	// Assert: Verify the response
	assert.NoError(t, err, "Ping should not return an error for valid domain")
	assert.Equal(t, http.StatusOK, statusCode, "Expected HTTP 200 status")
}

func TestPing_SuccessWithIP(t *testing.T) {
	// Arrange: Set up a mock server to simulate a successful ping
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "HEAD", r.Method, "Expected HEAD method")
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Extract the IP address from the mock server's URL
	host, port, _ := net.SplitHostPort(server.Listener.Addr().String())
	ip := net.JoinHostPort(host, port)

	// Act: Call Ping with the IP address
	statusCode, err := be.Ping(ip)

	// Assert: Verify the response
	assert.NoError(t, err, "Ping should not return an error for valid IP address")
	assert.Equal(t, http.StatusOK, statusCode, "Expected HTTP 200 status")
}

func TestPing_InvalidURL(t *testing.T) {
	// Act: Call Ping with an invalid target
	statusCode, err := be.Ping("invalid-url")

	// Assert: Verify the response
	assert.Error(t, err, "Ping should return an error for invalid URL")
	assert.Equal(t, 0, statusCode, "Status code should be 0 for invalid URL")
}

func TestPing_Failure(t *testing.T) {
	// Arrange: Set up a mock server to simulate a failure
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	// Act: Call Ping with the mock server's URL
	statusCode, err := be.Ping(server.URL[len("http://"):])

	// Assert: Verify the response
	assert.NoError(t, err, "Ping should not return a fatal error even on server failure")
	assert.Equal(t, http.StatusServiceUnavailable, statusCode, "Expected HTTP 503 status")
}

func TestPing_Timeout(t *testing.T) {
	// Arrange: Set up a mock server with a delayed response to simulate a timeout
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second) // Exceed the client timeout of 2 seconds
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Act: Call Ping with the mock server's URL
	statusCode, err := be.Ping(server.URL[len("http://"):])

	// Assert: Verify the timeout error
	assert.Error(t, err, "Ping should return an error on timeout")
	assert.Equal(t, 0, statusCode, "Status code should be 0 on timeout")
}
