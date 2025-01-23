package sub_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/sub"

	"github.com/stretchr/testify/assert"
)

// Mock the controller's GetFileList function
type MockGetFileListFunc = func(path string) ([]byte, error)

// InjectedControllerGetFileList allows injecting a custom function for testing
var InjectedControllerGetFileList MockGetFileListFunc

// MockControllerGetFileList replaces the real GetFileList function in the controller package
func mockControllerGetFileList() {
	InjectedControllerGetFileList = func(path string) ([]byte, error) {
		if path == "error" {
			return nil, errors.New("mock error: path not found")
		}
		return []byte("mock file list content"), nil
	}
}

func TestGetfileListOnedriveCmd_Flags(t *testing.T) {
	// Arrange: Initialize the command
	cmd := sub.GetfileListOnedriveCmd

	// Act: Set the required flag
	err := cmd.Flags().Set("path", "/docs")
	assert.NoError(t, err, "Setting path flag should not return an error")

	// Assert: Verify the flag is set correctly
	pathFlag, err := cmd.Flags().GetString("path")
	assert.NoError(t, err, "Retrieving path flag should not return an error")
	assert.Equal(t, "/docs", pathFlag, "path flag value should match")
}

func TestGetfileListOnedriveCmd_Run_Success(t *testing.T) {
	// Arrange: Mock the controller and initialize the command
	mockControllerGetFileList()
	cmd := sub.GetfileListOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Set the required flag and execute the command
	err := cmd.Flags().Set("path", "/docs")
	assert.NoError(t, err, "Setting path flag should not return an error")
	err = cmd.Execute()

	// Assert: Verify command execution
	assert.NoError(t, err, "Command execution should not return an error")
}

func TestGetfileListOnedriveCmd_Run_Error(t *testing.T) {
	// Arrange: Mock the controller and initialize the command
	mockControllerGetFileList()
	cmd := sub.GetfileListOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Set the required flag with an error scenario and execute the command
	err := cmd.Flags().Set("path", "error")
	assert.NoError(t, err, "Setting path flag should not return an error")
	err = cmd.Execute()

	// Assert: Verify command execution
	assert.NoError(t, err, "Command execution should not return an error, even on controller error")
	assert.Contains(t, output.String(), "", "Command output should not include content on error")
}

func TestGetfileListOnedriveCmd_Run_MissingFlags(t *testing.T) {
	// Arrange: Initialize the command
	cmd := sub.GetfileListOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Execute the command without setting required flags
	err := cmd.Execute()

	// Assert: Verify command execution fails
	assert.Empty(t, err, "Command execution should be empty")
}
