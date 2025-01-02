package azure

import (
	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive"
)

// AzureCmd represents the azure command
var AzureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Azure-related commands",
	Long:  "Commands to manage Azure services such as OneDrive, office365, azure resources, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help() // Show help if no subcommands are provided
	},
}

func init() {
	// Add the OneDrive command to Azure
	AzureCmd.AddCommand(oneDrive.OnedriveCmd)
}
