package controller

import (
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

func UploadFile(localPath string, pathToUpload string) ([]byte, error) {
	// Declare variables
	var data map[string]interface{}
	// Call the GetDriveItems function from the backend
	result, err := be.UploadItem(localPath, pathToUpload)
	// Create interface for data return
	data = map[string]interface{}{
		"result": result,
	}
	// Handle the response
	return helpers.HandleControllerGeneric(
		"File uploaded successfully!",
		"UploadFile",
		data,
		err,
	)
}
