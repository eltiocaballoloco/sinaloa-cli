package helpers

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/errors"
)

// HandleResponseError simplifies the creation and display of JSON-formatted error messages
func HandleResponseError(message, code string, err error) []byte {
	var errorMessage string

	if os.Getenv("SINALOA_CLI_DEBUG") == "" || os.Getenv("SINALOA_CLI_DEBUG") == "false" {
		errorMessage = fmt.Sprintf("%s", err.Error())
	} else {
		errorMessage = fmt.Sprintf("%s --> %s", message, err.Error())
	}

	errorResponse := errors.NewErrorResponse(false, code, errorMessage)

	errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
	if jsonErr != nil {
		errJSON := map[string]interface{}{
			"response": false,
			"code":     "500",
			"message":  fmt.Sprintf("Error marshaling JSON: %v", jsonErr),
			"data":     struct{}{},
		}

		errorJsonResponse, _ = json.MarshalIndent(errJSON, "", "  ")
		fmt.Println(string(errorJsonResponse))
		return errorJsonResponse
	}

	fmt.Println(string(errorJsonResponse))

	return errorJsonResponse
}
