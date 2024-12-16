package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/net/be"
	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
)

var client = &http.Client{
	Timeout: 2 * time.Second,
}

func Ping(urlPath string) ([]byte, error) {
	// Call the Ping function from the backend
	if statusCode, err := be.Ping(urlPath); err != nil {
		// Handle the error
		errorMessage := fmt.Sprintf("Ping error: %s", err.Error())
		return helpers.HandleController(
			false,
			"500",
			errorMessage,
			"Ping",
			struct{}{},
			err,
		)
	} else {
		// Handle the success response
		successData := struct {
			StatusCode int `json:"status_code"`
		}{
			StatusCode: statusCode,
		}
		// Handle the response
		return helpers.HandleController(
			true,
			"200",
			"Ping successful",
			"Ping",
			successData,
			err,
		)
	}
}
