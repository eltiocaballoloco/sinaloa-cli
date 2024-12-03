package oneDrive

import (
	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/azure/oneDrive/sub"
)

// OnedriveCmd represents the one-drive command
var OnedriveCmd = &cobra.Command{
	Use:   "one-drive",
	Short: "OneDrive-related commands",
	Long:  "Commands to interact with Microsoft OneDrive, such as fetching files, uploading, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // Show help if no subcommands are provided
	},
}

func init() {
	// Add subcommands to the OneDrive command
	OnedriveCmd.AddCommand(sub.GetfileOnedriveCmd)
	OnedriveCmd.AddCommand(sub.GetfileListOnedriveCmd)
}
