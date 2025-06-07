package helpers

import (
	"strings"
)

func ReturnCompleteDockerRepoPath(repoUrl string, dockerRepo string) string {
	trimmed := strings.TrimPrefix(repoUrl, "git@gitlab.com:")
	trimmed = strings.TrimSuffix(trimmed, ".git")
	dockerRepoPath := dockerRepo + "/" + trimmed
	return dockerRepoPath
}
