package net_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net"
)

func TestNetCmd_HasSubCommands(t *testing.T) {
	// Act: Get the list of subcommands
	subCommands := net.NetCmd.Commands()

	// Assert: Verify that PingCmd (or other subcommands) is added
	assert.NotNil(t, subCommands, "NetCmd should have subcommands")
	assert.True(t, len(subCommands) > 0, "NetCmd should have at least one subcommand")

	// Verify specific subcommands exist (e.g., PingCmd)
	foundPing := false
	for _, cmd := range subCommands {
		if cmd.Name() == "ping" {
			foundPing = true
			break
		}
	}
	assert.True(t, foundPing, "NetCmd should include the 'ping' subcommand")
}

func TestNetCmd_ExecutesHelp(t *testing.T) {
	// Arrange: Prepare a buffer to capture command output
	buf := new(bytes.Buffer)
	net.NetCmd.SetOut(buf)
	net.NetCmd.SetArgs([]string{}) // No arguments to trigger the Help message

	// Act: Execute the command
	err := net.NetCmd.Execute()

	// Assert: Verify that the Help message is displayed
	assert.NoError(t, err, "Executing NetCmd without arguments should not return an error")
	output := buf.String()
	assert.Contains(t, output, "Net is a palette that contains network-based commands", "Help message should be displayed")
}
