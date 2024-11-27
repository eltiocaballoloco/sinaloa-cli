package vault

import (
	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/vault/sub"
)

var VaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Short description of vault",
	Long:  "Long description of vault",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	VaultCmd.AddCommand(sub.SetsecretVaultCmd)
}
