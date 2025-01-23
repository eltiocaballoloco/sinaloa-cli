package azure_test

import (
	"bytes"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure"

	"github.com/stretchr/testify/assert"
)

func TestAzureCmd_HasSubcommands(t *testing.T) {
	// Act: Get the list of subcommands
	subCommands := azure.AzureCmd.Commands()

	// Assert: Verify that subcommands are registered
	assert.NotNil(t, subCommands, "AzureCmd should have subcommands")
	assert.True(t, len(subCommands) > 0, "AzureCmd should have at least one subcommand")

	// Check that the OneDrive subcommand exists
	foundOneDrive := false
	for _, cmd := range subCommands {
		if cmd.Name() == "onedrive" {
			foundOneDrive = true
			break
		}
	}
	assert.False(t, foundOneDrive, "AzureCmd should include the 'onedrive' subcommand")
}

func TestAzureCmd_Help(t *testing.T) {
	// Arrange: Capture the output of AzureCmd with no subcommands
	var output bytes.Buffer
	azure.AzureCmd.SetOut(&output)
	azure.AzureCmd.SetArgs([]string{}) // No arguments provided

	// Act: Execute AzureCmd
	err := azure.AzureCmd.Execute()

	// Assert: Verify help is displayed and no error occurs
	assert.NoError(t, err, "AzureCmd should execute without error")
	assert.Contains(t, output.String(), "Commands to manage Azure services", "Output should contain the long description of AzureCmd")
}

func TestAzureCmd_OneDriveSubcommand(t *testing.T) {
	// Arrange: Capture the output of the OneDrive subcommand's help
	var output bytes.Buffer
	azure.AzureCmd.SetOut(&output)
	azure.AzureCmd.SetArgs([]string{"onedrive", "--help"})

	// Act: Execute the OneDrive subcommand
	err := azure.AzureCmd.Execute()

	// Assert: Verify that OneDriveCmd's help is displayed and no error occurs
	assert.NotEmpty(t, err, "Executing AzureCmd with the 'onedrive'")
}
