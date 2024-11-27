package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/net/be"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
}

func Ping(urlPath string) ([]byte, error) {
	if statusCode, err := be.Ping(urlPath); err != nil {
		errorMessage := fmt.Sprintf("Ping error: %s", err.Error())
		errorResponse := errors.NewErrorResponse(false, "500", errorMessage)
		errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("Error marshaling JSON:", jsonErr)
		}
		return errorJsonResponse, err
	} else {
		successData := struct {
			StatusCode int `json:"status_code"`
		}{
			StatusCode: statusCode,
		}
		successResponse := response.NewResponse(true, "200", "Ping successful", successData)
		jsonResponse, jsonErr := json.MarshalIndent(successResponse, "", "  ")
		if jsonErr != nil {
			fmt.Println("Error marshaling JSON:", jsonErr)
		}
		return jsonResponse, err
	}
}
