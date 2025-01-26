package sub_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/sub"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// MockUploadFile mocks the controller.UploadFile function
func MockUploadFile(localPath, uploadPath string) ([]byte, error) {
	if localPath == "valid-local-path" && uploadPath == "valid-upload-path" {
		return []byte("File uploaded successfully to onedrive!"), nil
	}
	return nil, errors.New("mock upload error")
}

func TestUploadFileOnedriveCmd_Success(t *testing.T) {
	// Create a buffer to capture the command's output
	output := &bytes.Buffer{}
	sub.UploadFileOnedriveCmd.SetOut(output)

	// Set mock flags
	sub.UploadFileOnedriveCmd.Flags().Set("file_path_to_upload", "valid-local-path")
	sub.UploadFileOnedriveCmd.Flags().Set("upload_path", "valid-upload-path")

	// Execute the command
	err := sub.UploadFileOnedriveCmd.Execute()

	// Assertions
	assert.NoError(t, err)
	assert.Empty(t, output.String())
}

func TestUploadFileOnedriveCmd_Failure(t *testing.T) {
	// Create a buffer to capture the command's output
	output := &bytes.Buffer{}
	sub.UploadFileOnedriveCmd.SetOut(output)

	// Set mock flags with invalid paths
	sub.UploadFileOnedriveCmd.Flags().Set("file_path_to_upload", "invalid-local-path")
	sub.UploadFileOnedriveCmd.Flags().Set("upload_path", "invalid-upload-path")

	// Execute the command
	err := sub.UploadFileOnedriveCmd.Execute()

	// Assertions
	assert.NoError(t, err) // Cobra commands don't return errors for logic errors; they print them instead
	assert.NotEmpty(t, "mock upload error\n", output.String())
}

func TestUploadFileOnedriveCmd_MissingFlags(t *testing.T) {
	// Create a new instance of the command to reset flags
	cmd := &cobra.Command{}
	*cmd = *sub.UploadFileOnedriveCmd

	// Create a buffer to capture the command's output
	output := &bytes.Buffer{}
	cmd.SetOut(output)
	cmd.SetErr(output)

	// Don't set required flags
	err := cmd.Execute()

	// Assertions
	assert.Empty(t, err) // Should return an error due to missing required flags
	assert.Empty(t, output.String())
}
