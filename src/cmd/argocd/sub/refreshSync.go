package sub

import (
	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/controller"
	"github.com/eltiocaballoloco/sinaloa-cli/src/helpers"
)

var (
	gitId      string
	gitlabPath string
	env        string
	regions    string
)

var SyncArgocdCmd = &cobra.Command{
	Use:   "sync",
	Short: "ArgoCD sync apps",
	Long:  "Sync applications from ArgoCD",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration from .env
		helpers.LoadConfig()

		// Start the argocd sync
		controller.RefreshSync(
			gitId,
			gitlabPath,
			env,
			regions,
			helpers.AppConfig.ARGOCD_URL,
			helpers.AppConfig.ARGOCD_USER,
			helpers.AppConfig.ARGOCD_PASSWORD,
		)
	},
}

func init() {
	SyncArgocdCmd.Flags().StringVarP(&gitId, "git-id", "g", "", "Git id of the application")
	SyncArgocdCmd.Flags().StringVarP(&gitlabPath, "gitlab-path", "p", "", "Gitlab path of the application")
	SyncArgocdCmd.Flags().StringVarP(&env, "env", "e", "", "Environments like dev, test, prod")
	SyncArgocdCmd.Flags().StringVarP(&regions, "regions", "r", "", "Regions, if an application need to be deployed in multiple clusters at the same time")
}
