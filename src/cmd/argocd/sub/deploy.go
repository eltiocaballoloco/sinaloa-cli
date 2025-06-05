package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	json string
)

var DeployArgocdCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Short description of deploy",
	Long:  "Long description of deploy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing deploy in argocd")

		// 1 - Checking the image version
		//     * incremental: we need to get last version from the docker hub
		//                    and sostitute in the values
		//     * latest:      It is setup directly in the values file
		//     * void '':     Take the version from the values file

		// 2 - Take the secrets for env, from onedrive

		// 3 - Concat extra_secrets

		// 4 - if module is different from '', take the secret from onedrive
		//     based on the specific module and the environment

		// 5 - Concat chart params

		// 6 - Create iwth previuos points the helm template to render on stdout for argocd
	},
}

func init() {
	DeployArgocdCmd.Flags().StringVarP(&json, "json", "j", "", "Json to pass for the deploy with argocd")
}
