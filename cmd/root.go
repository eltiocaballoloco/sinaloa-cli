package cmd

import (
	"os"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/api"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/env"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/haproxy"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/net"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/storj"
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sinaloa",
	Short: "The sinaloa cli",
	Long:  `The sinaloa cli`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func addSubcommandPalettes() {
	rootCmd.AddCommand(haproxy.HaproxyCmd)
	rootCmd.AddCommand(storj.StorjCmd)
	rootCmd.AddCommand(env.EnvCmd)
	rootCmd.AddCommand(net.NetCmd)
	rootCmd.AddCommand(api.ApiCmd)
	rootCmd.AddCommand(version.VersionCmd)
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	addSubcommandPalettes()
}
