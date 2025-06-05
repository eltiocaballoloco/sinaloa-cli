package be

import (
	"fmt"
	"net/url"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/docker"
)

// DeleteImages deletes a list of tags from the given Docker Hub repository path
// Returns lists of tags successfully deleted and not deleted
func DeleteImages(token string, repoPath string, tags []docker.TagInfoInternal) (map[string]interface{}, error) {
	// Declare variables
	client := helpers.NewApiClient("https://hub.docker.com", token, "Bearer")
	tagsDeleted := []string{}
	tagsNotDeleted := []string{}

	for _, tag := range tags {
		// URL-encode tag to handle special characters safely
		encodedTag := url.PathEscape(tag.Name)

		// DELETE endpoint for tag: /v2/repositories/{repoPath}/tags/{tag}/
		endpoint := fmt.Sprintf("/v2/repositories/%s/tags/%s/", repoPath, encodedTag)
		resp := client.Request("DELETE", endpoint, nil)

		// Check if the response is successful
		// If the response is successful, append the tag to tagsDeleted
		// If the response is not successful, append the tag to tagsNotDeleted
		if resp.Response {
			tagsDeleted = append(tagsDeleted, tag.Name)
		} else {
			tagsNotDeleted = append(tagsNotDeleted, tag.Name)
		}
	}

	// Prepare the result map with lists of deleted and not deleted tags
	result := map[string]interface{}{
		"tags_deleted": map[string]interface{}{
			"tags_list": tagsDeleted,
			"count":     len(tagsDeleted),
		},
		"tags_not_deleted": map[string]interface{}{
			"tags_list": tagsNotDeleted,
			"count":     len(tagsNotDeleted),
		},
	}

	return result, nil
}
