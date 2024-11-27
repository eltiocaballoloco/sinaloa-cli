package storj

import (
	"github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub" // Import the sub package
	"github.com/spf13/cobra"
)

var StorjCmd = &cobra.Command{
	Use:   "storj",
	Short: "Short description of storj",
	Long:  "Long description of storj",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	StorjCmd.AddCommand(sub.AddStorjCmd)
	StorjCmd.AddCommand(sub.GetStorjCmd)
	StorjCmd.AddCommand(sub.DeleteStorjCmd)
}
