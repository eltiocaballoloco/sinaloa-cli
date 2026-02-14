package sub

import (
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/github/controller"
	"github.com/spf13/cobra"
)

var (
	getReposOrganization string
	getReposQuery        string
	getReposSaveJSON     bool
	getReposSavePathJSON string
)

var GetReposCmd = &cobra.Command{
	Use:   "get-repos",
	Short: "Fetch and filter GitHub repositories from an organization",
	Long: `Fetch and filter GitHub repositories from an organization.
	
Example:
  sinaloa github get-repos \
    --organization "OrgName" \
    --query "-my-repo-,repo-" \
    --save-json true \
    --save-path-json "/tmp/repos.json"`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the GetRepos controller
		result, err := controller.GetRepos(
			getReposOrganization,
			getReposQuery,
			getReposSaveJSON,
			getReposSavePathJSON,
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
	GetReposCmd.Flags().StringVarP(&getReposOrganization, "organization", "o", "", "Organization name (required)")
	GetReposCmd.Flags().StringVarP(&getReposQuery, "query", "q", "", "Comma-separated repo patterns to search")
	GetReposCmd.Flags().BoolVarP(&getReposSaveJSON, "save-json", "z", false, "Save to file (true) or display in terminal (false)")
	GetReposCmd.Flags().StringVarP(&getReposSavePathJSON, "save-path-json", "j", "", "Absolute path for JSON output file")

	GetReposCmd.MarkFlagRequired("organization")
}
