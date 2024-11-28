package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	file          string
	path_to_store string
)

var GetfileOnedriveCmd = &cobra.Command{
	Use:   "get-file",
	Short: "Short description of get-file",
	Long:  "Long description of get-file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing get-file in one-drive")
	},
}

func init() {
	GetfileOnedriveCmd.Flags().StringVarP(&file, "file", "f", "", "file to get from one drive")
	GetfileOnedriveCmd.MarkFlagRequired("file")
	GetfileOnedriveCmd.Flags().StringVarP(&path_to_store, "path_to_store", "p", "", "Path where you want store the file from one drive")
	GetfileOnedriveCmd.MarkFlagRequired("path_to_store")
}
