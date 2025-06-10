package shared

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/controller"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/docker"
)

func FetchLatestTag(repoUrl string, dockerRepo string) (string, error) {
	// Get the complete dockerhub path from repoUrl
	dockerRepoPath := helpers.ReturnCompleteDockerRepoPath(repoUrl, dockerRepo)

	// Get image list
	imageListBytes, err := controller.GetImages(
		dockerRepoPath,
		"100",
		"",
		"get",
		true,
	)
	if err != nil {
		return "error", fmt.Errorf("[Error] failed to fetch image list: %v", err)
	}

	var response docker.DockerHubResponse
	if err := json.Unmarshal(imageListBytes, &response); err != nil {
		return "error", fmt.Errorf("[Error] failed to parse image list JSON: %v", err)
	}

	// Get highest semver tag
	highest, err := getHighestSemverFromTags(response.Data.TagList)
	if err != nil {
		return "error", err
	}

	return highest, nil
}

func getHighestSemverFromTags(tags []docker.TagInfoInternal) (string, error) {
	semverRe := regexp.MustCompile(`^v?(\d+\.\d+\.\d+)$`)
	var versions []*semver.Version
	tagToOriginal := make(map[string]string)

	for _, tag := range tags {
		tagName := strings.TrimSpace(tag.Name)

		if !semverRe.MatchString(tagName) {
			fmt.Printf("Skipping non-matching tag: %s\n", tagName)
			continue
		}

		matches := semverRe.FindStringSubmatch(tagName)
		if len(matches) != 2 {
			fmt.Printf("Skipping malformed tag: %s\n", tagName)
			continue
		}

		normalized := matches[1] // e.g., 0.1.0
		v, err := semver.NewVersion(normalized)
		if err != nil {
			fmt.Printf("Skipping invalid semver tag: %s (%v)\n", tagName, err)
			continue
		}

		versions = append(versions, v)
		tagToOriginal[v.String()] = tagName
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("[Error] no valid semver tags found")
	}

	sort.Sort(semver.Collection(versions))
	highest := versions[len(versions)-1]

	return tagToOriginal[highest.String()], nil
}
