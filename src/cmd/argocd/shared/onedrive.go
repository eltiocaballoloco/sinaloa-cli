package shared

import (
	"fmt"

	"strings"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

func FetchSecret(
	env string,
	localPathToSaveFile string,
	repoUrl string,
	dockerRepo string,
	manifestToDownload string,
) error {
	// Return the path of the repo on docker (it is the same used across tools)
	secretPathOnOneDrive := helpers.ReturnCompleteDockerRepoPath(
		repoUrl,
		dockerRepo,
	)

	// Remove the prefix of the docker repo registry
	secretPathOnOneDrive = strings.TrimPrefix(secretPathOnOneDrive, dockerRepo+"/")

	// Get the manifest from the onedrive
	_, err := controller.GetFile(
		"development"+"/"+secretPathOnOneDrive+"/"+env+"/"+manifestToDownload,
		localPathToSaveFile,
	)
	if err != nil {
		return fmt.Errorf("[Error] Failed to fetch manifest from OneDrive (FetchSecret): %v", err)
	}

	return nil
}

func FetchExtraSecrets(
	env string,
	files string,
	repoUrl string,
	dockerRepo string,
	outputDir string,
) error {
	// Split the files string into a slice
	fileList := strings.Split(files, ",")

	// Loop through each file and fetch secrets
	for _, file := range fileList {
		errSecret := FetchSecret(
			env,
			outputDir+"/"+file,
			repoUrl,
			dockerRepo,
			file,
		)
		if errSecret != nil {
			return fmt.Errorf("[Error] Failed to fetch manifest from OneDrive (FetchExtraSecrets): %v", errSecret)
		}
	}

	return nil
}
