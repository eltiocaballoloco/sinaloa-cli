package cmd

import (
	"os"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/argocd"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/github"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net"
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sinaloa",
	Short: "The sinaloa cli",
	Long:  `The sinaloa cli`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true, // Disable the "completion" command
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPalettes() {
	rootCmd.AddCommand(argocd.ArgocdCmd)
	rootCmd.AddCommand(docker.DockerCmd)
	rootCmd.AddCommand(github.GithubCmd)
	rootCmd.AddCommand(azure.AzureCmd)
	rootCmd.AddCommand(net.NetCmd)
	rootCmd.AddCommand(version.VersionCmd)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubcommandPalettes()
}
