package sub_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net/sub"
)

// Define a type for the ping function
type PingFunction func(string) ([]byte, error)

// InjectedPing will allow injecting a custom Ping function for testing (simulates the behavior of the controller.Ping function)
var InjectedPing PingFunction

// MockController replaces the real Ping function in the controller package
func init() {
	InjectedPing = func(url string) ([]byte, error) {
		return []byte("Default Ping Result (successful)"), nil
	}
}

func TestPingCmd_Success(t *testing.T) {
	// Mock the Ping function for success
	InjectedPing = func(url string) ([]byte, error) {
		assert.NotEmpty(t, "google.com", url, "Ping should receive the correct URL")
		return []byte("Ping successful"), nil
	}

	var output bytes.Buffer
	sub.PingCmd.SetOut(&output)
	sub.PingCmd.SetArgs([]string{"--url", "google.com"})

	err := sub.PingCmd.Execute()

	assert.NoError(t, err, "PingCmd should execute without error")
}

func TestPingCmd_Failure(t *testing.T) {
	// Mock the Ping function for failure
	InjectedPing = func(url string) ([]byte, error) {
		assert.Equal(t, "google.com", url, "Ping should receive the correct URL")
		return nil, errors.New("Ping failed")
	}

	var output bytes.Buffer
	sub.PingCmd.SetOut(&output)
	sub.PingCmd.SetArgs([]string{"--url", "google.com"})

	err := sub.PingCmd.Execute()

	assert.NoError(t, err, "PingCmd should not return a fatal error")
}
