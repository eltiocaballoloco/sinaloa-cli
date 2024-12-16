package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
)

func HandleResponse(message string, code string, data interface{}) []byte {
	successResponse := response.NewResponse(true, code, message, data)
	successJsonResponse, err := json.MarshalIndent(successResponse, "", "  ")
	if err != nil {
		errorResponse := map[string]interface{}{
			"response": false,
			"code":     "500",
			"message":  fmt.Sprintf("Error marshaling JSON: %v", err),
			"data":     struct{}{},
		}

		errorJsonResponse, _ := json.MarshalIndent(errorResponse, "", "  ")
		fmt.Println(string(errorJsonResponse))
		return errorJsonResponse
	}

	fmt.Println(string(successJsonResponse))
	return successJsonResponse
}

func HandleController(result bool, statusCode string, message string, controllerFunction string, data interface{}, err error) ([]byte, error) {
	if result && err == nil {
		successResponse := response.NewResponse(result, statusCode, message, data)
		jsonResponse, jsonErr := json.MarshalIndent(successResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (new response):", jsonErr)
		}
		return jsonResponse, err
	} else {
		fmt.Printf("[Error] %v\n", err)
		errorResponse := errors.NewErrorResponse(result, statusCode, message)
		errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (error response):", jsonErr)
		}
		return errorJsonResponse, err
	}
}
