package argocd

import (
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/sub"

	"github.com/spf13/cobra"
)

var ArgocdCmd = &cobra.Command{
	Use:   "argocd",
	Short: "ArgoCD commands",
	Long:  "ArgoCD commands to manage ArgoCD deploy, applications, projects and other resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ArgocdCmd.AddCommand(sub.DeployArgocdCmd)
	ArgocdCmd.AddCommand(sub.SyncArgocdCmd)
}
