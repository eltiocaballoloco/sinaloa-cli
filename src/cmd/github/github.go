package github

import (
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/github/sub"
	"github.com/spf13/cobra"
)

var GithubCmd = &cobra.Command{
	Use:   "github",
	Short: "GitHub commands",
	Long:  "GitHub commands to manage repositories, analyze deployments, and extract deployment information",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	GithubCmd.AddCommand(sub.GetReposCmd)
	GithubCmd.AddCommand(sub.ReposDeployEnvironmentsCmd)
}

