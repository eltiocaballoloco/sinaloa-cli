package haproxy

import (
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/haproxy/sub"
	"github.com/spf13/cobra"
)

var HaproxyCmd = &cobra.Command{
	Use:   "haproxy",
	Short: "Short description of haproxy",
	Long:  "Long description of haproxy",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	HaproxyCmd.AddCommand(sub.ReceiveHaproxyCmd)
}
