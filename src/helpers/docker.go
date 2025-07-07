package helpers

import (
	"strings"
)

func ReturnCompleteDockerRepoPath(repoUrl string, dockerRepo string) string {
	trimmed := strings.TrimSuffix(repoUrl, "https://gitlab.com/")
	trimmed = strings.TrimSuffix(trimmed, ".git")
	dockerRepoPath := dockerRepo + "/" + strings.ReplaceAll(trimmed, "/", ".")
	return dockerRepoPath
}
