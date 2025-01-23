package controller

import (
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

func GetFileList(path string) ([]byte, error) {
	// Call the GetDriveItems function from the backend
	apiResponse, err := be.GetDriveItems(path)
	// Handle the response
	return helpers.HandleControllerApi(
		apiResponse.Response,
		fmt.Sprintf("%d", apiResponse.StatusCode),
		apiResponse.Message,
		"GetFileList",
		apiResponse.Body,
		err,
	)
}
