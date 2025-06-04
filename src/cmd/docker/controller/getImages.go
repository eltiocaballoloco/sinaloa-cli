package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/shared/auth"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/model/docker"
)

func GetImages(repoPath string, imagesForPage string) ([]byte, error) {
	// Get the token
	token, refresh, err := auth.LoginToDockerHub(
		helpers.AppConfig.DOCKER_HUB_USER_R,
		helpers.AppConfig.DOCKER_HUB_USER_R,
	)
	if err != nil {
		log.Fatalf("Docker auth error: %v", err)
	}

	// Get the tag list from BE layer
	result, statusCode, err := be.GetImages(token, refresh, repoPath, imagesForPage)
	if err != nil {
		errorMessage := fmt.Sprintf("Docker getImages error: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			strconv.Itoa(statusCode),
			errorMessage,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Unmarshal raw bytes into internal model
	var internalResp docker.TagResponseInternal
	if err := json.Unmarshal(result, &internalResp); err != nil {
		errorMessage := fmt.Sprintf("Failed to parse docker response: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			strconv.Itoa(statusCode),
			errorMessage,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Convert internal model into map[string]interface{}
	// Marshal struct to JSON bytes again
	internalJSON, err := json.Marshal(internalResp)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to marshal internal response from docker GetImages: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			strconv.Itoa(statusCode),
			errorMessage,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Unmarshal JSON bytes into map
	var dataMap map[string]interface{}
	if err := json.Unmarshal(internalJSON, &dataMap); err != nil {
		errorMessage := fmt.Sprintf("Failed to convert internal response to map from docker GetImages: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			strconv.Itoa(statusCode),
			errorMessage,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Return the response with data map
	return helpers.HandleControllerApi(
		true,
		strconv.Itoa(statusCode),
		"Docker images successfully retrieved",
		"GetImages",
		dataMap,
		nil,
	)
}
