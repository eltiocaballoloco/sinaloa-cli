package controller

import (
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/azure/oneDrive/be"
	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
)

func GetFileList(path string) ([]byte, error) {
	// Call the GetDriveItems function from the backend
	apiResponse, err := be.GetDriveItems(path)

	// Handle the response
	return helpers.HandleController(
		apiResponse.Response,
		fmt.Sprintf("%d", apiResponse.StatusCode),
		apiResponse.Message,
		"GetFileList",
		apiResponse.Body,
		err,
	)
}
