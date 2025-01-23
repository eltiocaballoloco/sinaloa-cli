package version_test

import (
	"bytes"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/version"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func executeCommand(cmd *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs(args)

	err := cmd.Execute()
	return buf.String(), err
}

func TestVersionCmdOutputNotNull(t *testing.T) {
	output, err := executeCommand(version.VersionCmd) // Access via import

	// Assertions
	assert.NoError(t, err, "Expected no error while executing version command")
	assert.Empty(t, output, "Output should be empty")
}
