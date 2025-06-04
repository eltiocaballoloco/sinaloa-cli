package controller

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

func GetImages(repoPath string, imagesForPage string, method string, dryRun bool) ([]byte, error) {
	// Get the token
	helpers.LoadConfig()
	token, refresh, err := shared.LoginToDockerHub(
		helpers.AppConfig.DOCKER_HUB_USER_R,
		helpers.AppConfig.DOCKER_HUB_PWD_R,
	)
	if err != nil {
		errorAuthFailed := fmt.Sprintf("Docker auth error: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			"500",
			errorAuthFailed,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Get the tag list from BE layer
	result, statusCode, err := be.GetImages(token, refresh, repoPath, imagesForPage)
	if err != nil {
		errorMessage := fmt.Sprintf("Docker GetImages error: %s", err.Error())
		return helpers.HandleControllerApi(
			false,
			strconv.Itoa(statusCode),
			errorMessage,
			"GetImages",
			struct{}{},
			err,
		)
	}

	// Method get images and dry-run true
	// (for get dry-run it is default to true)
	if method == "get" && dryRun {
		// Return the response with data map
		return helpers.HandleControllerApi(
			true,
			strconv.Itoa(statusCode),
			"Docker images successfully retrieved",
			"GetImages",
			result,
			nil,
		)
	}

	// Take the result and filter the tags to delete

	// Method delete images and dry-run true
	// (print only the tags to delete)
	if method == "delete" && dryRun {
		// Return the response with data map for dry run
		return helpers.HandleControllerApi(
			true,
			strconv.Itoa(statusCode),
			"Docker images to delete successfully retrieved",
			"DeleteImages",
			result,
			nil,
		)
	}

	// Method delete images and dry-run false
	// (delete the images)
	if method == "delete" && !dryRun {
		// Delete the images
		deletedImages, status, err := be.GetImages(token, refresh, repoPath, imagesForPage)
		if err != nil {
			errorMessage := fmt.Sprintf("Docker DeleteImages error: %s", err.Error())
			return helpers.HandleControllerApi(
				false,
				strconv.Itoa(status),
				errorMessage,
				"DeleteImages",
				struct{}{},
				err,
			)
		}

		return helpers.HandleControllerApi(
			true,
			strconv.Itoa(statusCode),
			"Docker images successfully deleted",
			"DeleteImages",
			deletedImages,
			nil,
		)
	}

	// If no method or dry-run specified, return an error
	return helpers.HandleControllerApi(
		false,
		"500",
		"No method or dry-run specified for Docker images",
		"GetImages-DeleteImages",
		struct{}{},
		errors.New("No method or dry-run specified for Docker images (GetImages-DeleteImages)"),
	)
}
