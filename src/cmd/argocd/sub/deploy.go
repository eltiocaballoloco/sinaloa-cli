package sub

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd/controller"
	"github.com/eltiocaballoloco/sinaloa-cli/src/models/argocd"
)

var (
	jsonInput string
)

var DeployArgocdCmd = &cobra.Command{
	Use:   "deploy",
	Short: "ArgoCD deploy cmd using the plugin",
	Long:  "Deploy an application using ArgoCD with the specified JSON configuration, through the argo-plugin cmp.",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the json input is provided
		if jsonInput == "" {
			fmt.Fprintln(os.Stderr, "[Error] JSON input is required (-j or --json)")
			os.Exit(1)
		}

		// Convert the json input string
		// with the associated class model
		var params argocd.ArgoCDDeployParams
		errParseClass := json.Unmarshal([]byte(jsonInput), &params)
		if errParseClass != nil {
			fmt.Fprintln(os.Stderr, "[Error] Failed to parse JSON input:", errParseClass)
			os.Exit(1)
		}

		// Execute the deploy
		errDeploy := controller.Deploy(params)
		if errDeploy != nil {
			fmt.Fprintln(os.Stderr, "[Error] Failed to deploy with ArgoCD... ", errDeploy)
		}
	},
}

func init() {
	DeployArgocdCmd.Flags().StringVarP(&jsonInput, "json", "j", "", "Json to pass for the deploy with argocd")
}
