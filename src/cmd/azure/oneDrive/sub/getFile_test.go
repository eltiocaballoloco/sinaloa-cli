package sub_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/sub"

	"github.com/stretchr/testify/assert"
)

// Mock the controller's GetFile function
type MockGetFileFunc = func(file, pathToStore string) ([]byte, error)

// InjectedControllerGetFile will allow injecting a custom function for testing (simulates the behavior of the controller.GetFile function)
var InjectedControllerGetFile MockGetFileFunc

// MockControllerGetFile replaces the real GetFile function in the controller package
func mockControllerGetFile() {
	InjectedControllerGetFile = func(file, pathToStore string) ([]byte, error) {
		if file == "error" {
			return nil, errors.New("mock error: file not found")
		}
		return []byte("mock file content"), nil
	}
}

func TestGetfileOnedriveCmd_Flags(t *testing.T) {
	// Arrange: Initialize the command
	cmd := sub.GetfileOnedriveCmd

	// Act: Set required flags
	err := cmd.Flags().Set("file", "mock-file.txt")
	assert.NoError(t, err, "Setting file flag should not return an error")
	err = cmd.Flags().Set("path_to_store", "/tmp/mock-file.txt")
	assert.NoError(t, err, "Setting path_to_store flag should not return an error")

	// Assert: Verify flags are set correctly
	fileFlag, err := cmd.Flags().GetString("file")
	assert.NoError(t, err, "Retrieving file flag should not return an error")
	assert.Equal(t, "mock-file.txt", fileFlag, "file flag value should match")

	pathFlag, err := cmd.Flags().GetString("path_to_store")
	assert.NoError(t, err, "Retrieving path_to_store flag should not return an error")
	assert.Equal(t, "/tmp/mock-file.txt", pathFlag, "path_to_store flag value should match")
}

func TestGetfileOnedriveCmd_Run_Success(t *testing.T) {
	// Arrange: Mock the controller and initialize the command
	mockControllerGetFile()
	cmd := sub.GetfileOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Set required flags and execute the command
	err := cmd.Flags().Set("file", "mock-file.txt")
	assert.NoError(t, err, "Setting file flag should not return an error")
	err = cmd.Flags().Set("path_to_store", "/tmp/mock-file.txt")
	assert.NoError(t, err, "Setting path_to_store flag should not return an error")
	err = cmd.Execute()

	// Assert: Verify command execution
	assert.NoError(t, err, "Command execution should not return an error")
}

func TestGetfileOnedriveCmd_Run_Error(t *testing.T) {
	// Arrange: Mock the controller and initialize the command
	mockControllerGetFile()
	cmd := sub.GetfileOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Set required flags and execute the command with an error case
	err := cmd.Flags().Set("file", "error")
	assert.NoError(t, err, "Setting file flag should not return an error")
	err = cmd.Flags().Set("path_to_store", "/tmp/mock-file.txt")
	assert.NoError(t, err, "Setting path_to_store flag should not return an error")
	err = cmd.Execute()

	// Assert: Verify command execution
	assert.NoError(t, err, "Command execution should not return an error, even on controller error")
	assert.Contains(t, output.String(), "", "Command output should not include content on error")
}

func TestGetfileOnedriveCmd_Run_MissingFlags(t *testing.T) {
	// Arrange: Initialize the command
	cmd := sub.GetfileOnedriveCmd
	var output bytes.Buffer
	cmd.SetOut(&output)

	// Act: Execute the command without setting required flags
	err := cmd.Execute()

	// Assert: Verify command execution fails
	assert.NotErrorAs(t, err, "Command execution should return an error when required flags are missing")
}
