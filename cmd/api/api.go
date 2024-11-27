package api

import (
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/api/sub" // Import the sub package
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var ApiCmd = &cobra.Command{
	Use:   "api",
	Short: "This command is used to start and manage the api server",
	Long:  `This command is used to start and manage the api server. Starting api server all commands of the cli will be available via api`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	ApiCmd.AddCommand(sub.StartApiCmd) // Use the StartApiCmd from the sub package
}
