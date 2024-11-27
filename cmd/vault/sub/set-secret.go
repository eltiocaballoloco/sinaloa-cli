package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	path string
)

var SetsecretVaultCmd = &cobra.Command{
	Use:   "set-secret",
	Short: "Short description of set-secret",
	Long:  "Long description of set-secret",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing set-secret in vault")
	},
}

func init() {
	SetsecretVaultCmd.Flags().StringVarP(&path, "path", "p", "", "path to upload the secret")
	SetsecretVaultCmd.MarkFlagRequired("path")
}
