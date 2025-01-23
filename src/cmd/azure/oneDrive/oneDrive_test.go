package oneDrive_test

import (
	"bytes"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/sub"

	"github.com/stretchr/testify/assert"
)

func TestOnedriveCmd_HasSubcommands(t *testing.T) {
	// Act: Get the list of subcommands
	subCommands := oneDrive.OnedriveCmd.Commands()

	// Assert: Verify subcommands are registered
	assert.NotNil(t, subCommands, "OnedriveCmd should have subcommands")
	assert.True(t, len(subCommands) > 0, "OnedriveCmd should have at least one subcommand")

	// Check that specific subcommands exist
	expectedSubCommands := []string{"get-file", "get-file-list"}
	for _, expected := range expectedSubCommands {
		found := false
		for _, cmd := range subCommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "OnedriveCmd should include the '%s' subcommand", expected)
	}
}

func TestOnedriveCmd_Help(t *testing.T) {
	// Arrange: Capture the output of OnedriveCmd with no subcommands
	var output bytes.Buffer
	oneDrive.OnedriveCmd.SetOut(&output)
	oneDrive.OnedriveCmd.SetArgs([]string{}) // No arguments provided

	// Act: Execute OnedriveCmd
	err := oneDrive.OnedriveCmd.Execute()

	// Assert: Verify help is displayed and no error occurs
	assert.NoError(t, err, "OnedriveCmd should execute without error")
	assert.Contains(t, output.String(), "Commands to interact with Microsoft OneDrive", "Output should contain the long description of OnedriveCmd")
}

func TestOnedriveCmd_SubcommandHelp(t *testing.T) {
	// Arrange: Capture the output of the subcommand help
	var output bytes.Buffer
	oneDrive.OnedriveCmd.SetOut(&output)
	oneDrive.OnedriveCmd.SetArgs([]string{"get-file", "--help"}) // Test help for a specific subcommand

	// Act: Execute the subcommand
	err := oneDrive.OnedriveCmd.Execute()

	// Assert: Verify the subcommand's help is displayed
	assert.NoError(t, err, "Executing OnedriveCmd with 'get-file --help' should not return an error")
	assert.Contains(t, output.String(), sub.GetfileOnedriveCmd.Short, "Output should contain the short description of the 'get-file' subcommand")
}

func TestOnedriveCmd_UnknownCommand(t *testing.T) {
	// Arrange: Capture the output for an unknown command
	var output bytes.Buffer
	oneDrive.OnedriveCmd.SetOut(&output)
	oneDrive.OnedriveCmd.SetArgs([]string{"unknown-command"}) // Invalid subcommand

	// Act: Execute the invalid command
	err := oneDrive.OnedriveCmd.Execute()

	// Assert: Verify error and output message
	assert.Error(t, err, "OnedriveCmd should return an error for unknown commands")
}
