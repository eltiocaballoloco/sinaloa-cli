package sub

import (
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/github/controller"
	"github.com/spf13/cobra"
)

var (
	deployEnvsOrganization    string
	deployEnvsQuery           string
	deployEnvsEnvs            string
	deployEnvsManifestSection string
	deployEnvsManifestName    string
	deployEnvsFolders         string
	deployEnvsSaveJSON        bool
	deployEnvsSavePathJSON    string
)

var ReposDeployEnvironmentsCmd = &cobra.Command{
	Use:   "repos-deploy-environments",
	Short: "Analyze deployment environments across multiple repositories and projects",
	Long: `Analyze deployment environments across multiple repositories and projects.
Scans repositories, parses manifest files, and generates a comprehensive deployment matrix.

Example:
  sinaloa github repos-deploy-environments \
    --organization "OrgName" \
    --query "my-repo-,repo-" \
    --envs "prod-,qa-,test-" \
    --folders "apps,src" \
    --save-json true \
    --save-path-json "/tmp/deployment-matrix.json"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the ReposDeployEnvironments controller
		result, err := controller.ReposDeployEnvironments(
			deployEnvsOrganization,
			deployEnvsQuery,
			deployEnvsEnvs,
			deployEnvsManifestSection,
			deployEnvsManifestName,
			deployEnvsFolders,
			deployEnvsSaveJSON,
			deployEnvsSavePathJSON,
		)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Print the result
		fmt.Println(string(result))
	},
}

func init() {
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsOrganization, "organization", "o", "", "Organization name (required)")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsQuery, "query", "q", "", "Comma-separated repo patterns to search")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsEnvs, "envs", "e", "prod-,qa-,test-", "Environment prefixes to search (comma-separated)")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsManifestSection, "manifest-section", "s", "environments", "YAML section to search")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsManifestName, "manifest-name", "m", "manifest.yaml", "Manifest filename")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsFolders, "folders", "f", "apps,src,micro-frontends", "Project folders to scan (comma-separated)")
	ReposDeployEnvironmentsCmd.Flags().BoolVarP(&deployEnvsSaveJSON, "save-json", "z", false, "Save to file (true) or display in terminal (false)")
	ReposDeployEnvironmentsCmd.Flags().StringVarP(&deployEnvsSavePathJSON, "save-path-json", "j", "", "Absolute path for JSON output file")

	ReposDeployEnvironmentsCmd.MarkFlagRequired("organization")
}
