package be

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/model/docker"
)

func GetImages(token string, refreshToken string, repoPath string, imagesForPage string) (docker.TagResponseInternal, int, error) {
	baseURL := "https://hub.docker.com"
	url := fmt.Sprintf("/v2/repositories/%s/tags?page_size=%s", repoPath, imagesForPage)
	client := helpers.NewApiClient(baseURL, token, "Bearer")

	var allResults []docker.TagResult
	var statusCode int

	for url != "" {
		response := client.Request("GET", url, nil)
		statusCode = response.StatusCode

		if !response.Response {
			return nil, statusCode, fmt.Errorf("failed to get images from dockerhub: %s", response.Message)
		}

		var tagsResp docker.TagsResponse
		if err := json.Unmarshal(response.Body, &tagsResp); err != nil {
			return nil, statusCode, fmt.Errorf("failed to unmarshal tags response: %w", err)
		}

		allResults = append(allResults, tagsResp.Results...)

		if tagsResp.Next != nil {
			url = *tagsResp.Next
			// strip full URL to relative path if needed
			url = url[len(baseURL):]
		} else {
			url = ""
		}
	}

	// Build final response
	final := docker.TagsResponse{
		Count:    len(allResults),
		Next:     nil,
		Previous: nil,
		Results:  allResults,
	}

	output, err := json.Marshal(final)
	if err != nil {
		return nil, statusCode, fmt.Errorf("failed to marshal final tags response: %w", err)
	}

	return convertTagsResponse(output), statusCode, nil
}

// convertTagsResponse converts the external TagsResponse into your internal TagResponseInternal.
func convertTagsResponse(external docker.TagsResponse) docker.TagResponseInternal {
	internalTags := make([]docker.TagInfoInternal, 0, len(external.Results))

	for _, tagResult := range external.Results {
		// For each external tag, convert to internal TagInfoInternal
		internalTag := docker.TagInfoInternal{
			Name:                tagResult.Name,
			IDImage:             tagResult.ID,
			IDRepository:        tagResult.Repository,
			IDCreator:           tagResult.Creator,
			LastUpdaterUsername: tagResult.LastUpdaterUsername,
			LastUpdated:         tagResult.LastUpdated,
			TagLastPulled:       tagResult.TagLastPulled,
			TagLastPushed:       tagResult.TagLastPushed,
			Digest:              tagResult.Digest,
		}
		internalTags = append(internalTags, internalTag)
	}

	return docker.TagResponseInternal{
		TagList: internalTags,
		Count:   len(internalTags),
	}
}
