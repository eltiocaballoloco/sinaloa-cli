package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCmd_Execute(t *testing.T) {
	// Arrange: Capture the output of the root command
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetArgs([]string{}) // No arguments

	// Act: Execute the root command
	err := rootCmd.Execute()

	// Assert: Verify that the command executed without error
	assert.NoError(t, err, "Root command should execute without error")
	assert.Contains(t, output.String(), "The sinaloa cli", "Output should contain the short description")
}

func TestRootCmd_HasSubcommands(t *testing.T) {
	// Arrange: Get the list of subcommands
	subCommands := rootCmd.Commands()

	// Assert: Verify that subcommands are registered
	assert.NotNil(t, subCommands, "Root command should have subcommands")
	assert.True(t, len(subCommands) > 0, "Root command should have at least one subcommand")

	// Verify specific subcommands exist
	expectedSubCommands := []string{"azure", "net", "version"}
	for _, expected := range expectedSubCommands {
		found := false
		for _, cmd := range subCommands {
			if cmd.Name() == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Root command should include the '%s' subcommand", expected)
	}
}

func TestRootCmd_HelpFlag(t *testing.T) {
	// Arrange: Capture the output of the --help flag
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetArgs([]string{"--help"})

	// Act: Execute the root command with --help
	err := rootCmd.Execute()

	// Assert: Verify that help text is displayed
	assert.NoError(t, err, "Root command with --help should execute without error")
	assert.Contains(t, output.String(), "The sinaloa cli", "Help output should contain the short description")
}

func TestRootCmd_ToggleFlag(t *testing.T) {
	// Arrange: Capture the output when using the toggle flag
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetArgs([]string{"--toggle"})

	// Act: Execute the root command with the toggle flag
	err := rootCmd.Execute()

	// Assert: Verify no error occurs and toggle flag is acknowledged
	assert.NoError(t, err, "Root command with --toggle should execute without error")
	assert.Contains(t, output.String(), "Help message for toggle", "Output should mention the toggle flag help")
}
