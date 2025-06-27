package controller

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/be"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/docker"
)

func GetImages(repoPath string, imagesForPage string, imagesToTake string, method string, dryRun bool) ([]byte, error) {
	// Get the token
	helpers.LoadConfig()
	token, refresh, err := shared.LoginToDockerHub(
		helpers.AppConfig.DOCKER_HUB_USER_RWD,
		helpers.AppConfig.DOCKER_HUB_PWD_RWD,
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
		// Return the response with data map for get
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
	imagesToTakeInt, err := strconv.Atoi(imagesToTake)
	if err != nil {
		// Manage the error
		imagesToTakeInt = 0 // Default value
		return helpers.HandleControllerApi(
			false,
			"500",
			err.Error(),
			"DeleteImages",
			struct{}{},
			err,
		)
	}

	// Tke the tags to delete
	resultDelete := filterTagsToDelete(result, imagesToTakeInt)

	// Method delete images and dry-run true
	// (print only the tags to delete)
	if method == "delete" && dryRun {
		// Return the response with data map for dry run
		return helpers.HandleControllerApi(
			true,
			strconv.Itoa(statusCode),
			"Docker images to delete successfully retrieved",
			"DeleteImages",
			resultDelete,
			nil,
		)
	}

	// Method delete images and dry-run false
	// (delete the images)
	if method == "delete" && !dryRun {
		// Delete the images
		deletedImages, err := be.DeleteImages(token, repoPath, resultDelete.TagList)

		// Return error if something fails
		if err != nil {
			errorMessage := fmt.Sprintf("Docker DeleteImages error: %s", err.Error())
			return helpers.HandleControllerApi(
				false,
				"500",
				errorMessage,
				"DeleteImages",
				struct{}{},
				err,
			)
		}

		// Return the response with data map for delete
		return helpers.HandleControllerApi(
			true,
			"200",
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

// filterTagsToDelete filters out the tags to delete from the original list
func filterTagsToDelete(result docker.TagResponseInternal, imagesToTake int) docker.TagResponseInternal {
	// Declare variables
	tagsToKeep := map[string]bool{
		"latest":   true,
		"unstable": true,
	}
	var otherTags []docker.TagInfoInternal

	// Look for the latest tag and separate other tags
	for _, tag := range result.TagList {
		if tag.Name == "latest" || tag.Name == "unstable" {
			tagsToKeep[tag.Name] = true
		} else {
			otherTags = append(otherTags, tag)
		}
	}

	// Filter semantic version tags and non-semantic version tags
	var semverTags []docker.TagInfoInternal
	var nonSemverTags []docker.TagInfoInternal
	for _, tag := range otherTags {
		versionStr := strings.TrimPrefix(tag.Name, "v")
		if _, err := semver.NewVersion(versionStr); err == nil {
			semverTags = append(semverTags, tag)
		} else {
			nonSemverTags = append(nonSemverTags, tag)
		}
	}

	// Order the semantic version tags in descending order
	sort.SliceStable(semverTags, func(i, j int) bool {
		vi, _ := semver.NewVersion(strings.TrimPrefix(semverTags[i].Name, "v"))
		vj, _ := semver.NewVersion(strings.TrimPrefix(semverTags[j].Name, "v"))
		return vi.GreaterThan(vj)
	})

	// Take first N semantic version tags
	for i := 0; i < imagesToTake && i < len(semverTags); i++ {
		tagsToKeep[semverTags[i].Name] = true
	}

	// Determinate tags to delete
	var tagsToDelete []docker.TagInfoInternal
	for _, tag := range result.TagList {
		if !tagsToKeep[tag.Name] {
			tagsToDelete = append(tagsToDelete, tag)
		}
	}

	// Return the filtered tags to delete
	return docker.TagResponseInternal{
		TagList: tagsToDelete,
		Count:   len(tagsToDelete),
	}
}
