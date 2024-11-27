package helpers

import (
	"encoding/json"
	"fmt"

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
