package cmd

import (
	"os"

	"github.com/spf13/cobra"
    "github.com/eltiocaballoloco/sinaloa-cli/cmd/vault"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/net"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/version"
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
    rootCmd.AddCommand(vault.VaultCmd)
	rootCmd.AddCommand(net.NetCmd)
	rootCmd.AddCommand(version.VersionCmd)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubcommandPalettes()
}
