package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/messages/response"
)

func HandleControllerApi(result bool, statusCode string, message string, controllerFunction string, data interface{}, err error) ([]byte, error) {
	// Convert data to a structured object if it is a byte slice
	if byteData, ok := data.([]byte); ok {
		// If data is []byte, convert it to a Go map or structured object for attachment
		var jsonData interface{}
		err := json.Unmarshal(byteData, &jsonData)
		if err != nil {
			fmt.Printf("[Error] Failed to unmarshal data: %v\n", err)
		} else {
			data = jsonData // Update data to hold parsed JSON
		}
	}
	// Check if the controller function was successful and
	// there are no errors in the conversion
	if result && err == nil {
		// Attach the parsed JSON data to the response
		successResponse := response.NewResponse(result, statusCode, message, data)
		// Marshal the response to JSON
		jsonResponse, jsonErr := json.MarshalIndent(successResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (new response):", jsonErr)
		}
		return jsonResponse, err
	} else {
		// Print an error message if the controller function failed
		fmt.Printf("[Error] An error occurred in the controller %s: %v\n", controllerFunction, err)
		errorResponse := errors.NewErrorResponse(result, statusCode, message)
		errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (error response):", jsonErr)
		}
		return errorJsonResponse, err
	}
}

// It has used to manage the response of the controller function
// if response return values instead of an api response
func HandleControllerGeneric(message string, controllerFunction string, data interface{}, err error) ([]byte, error) {
	// Convert data to a structured object if it is a byte slice
	if byteData, ok := data.([]byte); ok {
		// If data is []byte, convert it to a Go map or structured object for attachment
		var jsonData interface{}
		err := json.Unmarshal(byteData, &jsonData)
		if err != nil {
			fmt.Printf("[Error] Failed to unmarshal data: %v\n", err)
		} else {
			data = jsonData // Update data to hold parsed JSON
		}
	}
	// Check if the controller function was successful and
	// there are no errors in the conversion
	if data != nil && err == nil {
		// Attach the parsed JSON data to the response
		successResponse := response.NewResponse(true, "200", message, data)
		// Marshal the response to JSON
		jsonResponse, jsonErr := json.MarshalIndent(successResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (new response):", jsonErr)
		}
		return jsonResponse, err
	} else {
		// Print an error message if the controller function failed
		fmt.Printf("[Error] An error occurred in the controller %s: %v\n", controllerFunction, err)
		errorResponse := errors.NewErrorResponse(false, "500", "Error executing the command")
		errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("[Error] Controller", controllerFunction, ", error marshaling JSON (error response):", jsonErr)
		}
		return errorJsonResponse, err
	}
}
