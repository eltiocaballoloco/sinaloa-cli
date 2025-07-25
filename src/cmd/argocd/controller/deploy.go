package controller

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/shared"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/argocd"
)

func Deploy(params argocd.ArgoCDDeployParams) error {
	// 1 - Remove and create the tmp folder and copy the values.yaml

	// 2 - Checking the image version
	//     * incremental: we need to get last version from the docker hub
	//                    and sostitute in the values.
	//     * latest:      Use always the latest image instead of a specific version.
	//     * unstable:    For develop.
	//     * void '':     Take the version from the values file.

	// 3 - Take the secrets for env, from onedrive

	// 4 - Take extra_secrets

	// 5 - if module is different from '', take the secret from onedrive
	//     based on the specific module and the environment

	// 6 - If isIncremental/isLatest is true, replace the last image tag
	//     version in the values.yaml file

	// 7 - Create with previuos points the helm template
	//     to render on stdout for argocd

	// ----------------------------------------------------------------------------------

	// Declare variables and fix the work directory
	var outputDir string = filepath.Join("/tmp", params.AppName)
	var valuesYml string = "values.yaml"
	var secretYml string = "secret.yaml"
	var moduleValuesYml string = ""
	var moduleSecretYml string = ""

	// Step 1 - Remove and create the tmp folder and copy the values.yaml
	errRmvTmpFolder := os.RemoveAll(outputDir) // Remove the old tmp folder if exists
	if errRmvTmpFolder != nil {
		fmt.Fprintln(os.Stderr, "[Error] The deletion of the old tmp folder got an error in ArgoCD.Deploy:", errRmvTmpFolder)
		return errRmvTmpFolder
	}

	// Create the output directory
	errCreateTmpFolder := os.MkdirAll(outputDir, 0755)
	if errCreateTmpFolder != nil {
		fmt.Fprintln(os.Stderr, "[Error] The creation of the tmp folder got an error in ArgoCD.Deploy:", errCreateTmpFolder)
		return errCreateTmpFolder
	}

	// Copy the values.yaml file to the tmp folder
	errCopyValues := helpers.CopyFile(valuesYml, outputDir+"/"+valuesYml)
	if errCopyValues != nil {
		fmt.Fprintln(os.Stderr, "[Error] The copy of the values in the tmp folder got an error in ArgoCD.Deploy:", errCopyValues)
		return errCopyValues
	}

	// Step 2 - Image version
	var imageTag string
	var errDockerHub error
	var isIncremental bool = false
	var isLatest bool = false
	var isUnstable bool = false
	switch strings.ToLower(params.Tag) {
	case "incremental":
		// Get highest version if set to incremental
		imageTag, errDockerHub = shared.FetchLatestTag(params.RepoURL, params.DockerRepo)
		if errDockerHub != nil {
			fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.FetchLatestTag:", errDockerHub)
			return errDockerHub
		}
		// Set isIncremental to true, to replace after
		// the highest image tag inside the values.yaml
		isIncremental = true
	case "latest":
		// In this case, will use the image tag provided in the values.yaml
		imageTag = "latest"
		isLatest = true
	case "unstable":
		imageTag = "unstable"
		isUnstable = true
	default:
		// In this case, will use the image tag provided in the values.yaml
		imageTag = params.Tag
	}

	// Step 3 - Fetch base secret
	errSecret := shared.FetchSecret(
		params.Profile,
		outputDir+"/"+secretYml,
		params.RepoURL,
		params.DockerRepo,
		secretYml,
	)
	if errSecret != nil {
		fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.FetchSecret:", errSecret)
		return errSecret
	}

	// Step 4 - Add extra secrets (if defined)
	if params.ExtraSecrets != "" {
		// Get extra secrets from one drive (if configured)
		errExtraSecret := shared.FetchExtraSecrets(
			params.Profile,
			params.ExtraSecrets,
			params.RepoURL,
			params.DockerRepo,
			outputDir,
		)
		if errExtraSecret != nil {
			fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.FetchExtraSecrets:", errExtraSecret)
			return errExtraSecret
		}
	}

	// Step 5 - Module-based secret (if module is defined)
	if params.Module != "" {
		// Set the name of the manifests for the module
		moduleValuesYml = fmt.Sprintf("values-%s.yaml", params.Module)
		moduleSecretYml = fmt.Sprintf("secret-%s.yaml", params.Module)

		// Execute the fetch
		errModuleSecret := shared.FetchSecret(
			params.Profile,
			outputDir+"/"+moduleSecretYml,
			params.RepoURL,
			params.DockerRepo,
			moduleSecretYml,
		)
		if errModuleSecret != nil {
			fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.FetchSecretModule:", errModuleSecret)
			return errModuleSecret
		}
	}

	// Step 6 - Replace last image tag version in
	// the values if isIncremental or isLatest or isUnstable == true
	if isIncremental || isLatest || isUnstable {
		// Replace the image tag in the values.yaml file
		errImageV := helpers.UpdateImageTagWithRegex(
			filepath.Join(outputDir, valuesYml),
			"\""+imageTag+"\"",
		)
		if errImageV != nil {
			fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.UpdateImageTagWithRegex:", errImageV)
			return errImageV
		}
	}

	// Step 7 - Run helm template cmd
	args := []string{ // Create args string array
		"template",
		"--release-name", params.ReleaseName,
		"--namespace", params.Namespace,
		params.ChartName,
		"--repo", params.ChartRepo,
		"-f", outputDir + "/" + valuesYml,
		"-f", outputDir + "/" + secretYml,
	}

	// If there are extra values secrets files, append them to the cmd
	if params.ExtraSecrets != "" {
		fileList := strings.Split(params.ExtraSecrets, ",")
		for _, file := range fileList {
			args = append(args, "-f", outputDir+"/"+file)
		}
	}

	// If module is defined, append module-specific
	// values and secret to the cmd
	if params.Module != "" {
		args = append(args, "-f", outputDir+"/"+moduleValuesYml)
		args = append(args, "-f", outputDir+"/"+moduleSecretYml)
	}

	// If chartsParams contains something append them to the cmd
	if params.ChartParams != "" {
		args = append(args, params.ChartParams)
	}

	// Execute the Helm command and print real-time output to stdout/stderr
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	errCmd := cmd.Run()
	if errCmd != nil {
		fmt.Fprintln(os.Stderr, "[Error] ArgoCD.Deploy.Cmd.Run:", errCmd)
		return errCmd
	}

	return nil
}
