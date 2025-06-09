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
	fmt.Printf("[Info] Fetching latest tag from Docker repository: %s\n", dockerRepoPath)

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

	var response docker.TagResponseInternal
	if err := json.Unmarshal(imageListBytes, &response); err != nil {
		return "error", fmt.Errorf("[Error] failed to parse image list JSON: %v", err)
	}

	// Get highest semver tag
	highest, err := getHighestSemverFromTags(response.TagList)
	if err != nil {
		return "error", err
	}

	return highest, nil
}

func getHighestSemverFromTags(tags []docker.TagInfoInternal) (string, error) {
	semverRe := regexp.MustCompile(`^v?(\d+\.\d+\.\d+)$`)
	var versions []*semver.Version
	tagToOriginal := make(map[string]string)

	// Detect prefix convention
	prefixedWithV := false
	for _, tag := range tags {
		if strings.HasPrefix(tag.Name, "v") && semverRe.MatchString(tag.Name) {
			prefixedWithV = true
			break
		}
	}

	// Normalize and collect valid versions
	for _, tag := range tags {
		matches := semverRe.FindStringSubmatch(tag.Name)
		if len(matches) > 1 {
			normalized := matches[1]
			v, err := semver.NewVersion(normalized)
			if err == nil {
				versions = append(versions, v)
				// Map to original tag: if `v` is used, restore it
				if prefixedWithV {
					tagToOriginal[v.String()] = "v" + v.String()
				} else {
					tagToOriginal[v.String()] = v.String()
				}
			}
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("[Error] no valid semver tags found")
	}

	sort.Sort(semver.Collection(versions))
	highest := versions[len(versions)-1]

	return tagToOriginal[highest.String()], nil
}
