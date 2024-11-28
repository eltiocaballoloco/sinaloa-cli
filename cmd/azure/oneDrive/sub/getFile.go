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
	Short: "Get a file from one drive",
	Long:  "Get a file from one drive",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Executing get-file in one-drive")
	},
}

func init() {
	GetfileOnedriveCmd.Flags().StringVarP(&file, "file", "f", "", "file to get from one drive, example: -f secrets/file.txt")
	GetfileOnedriveCmd.MarkFlagRequired("file")
	GetfileOnedriveCmd.Flags().StringVarP(&path_to_store, "path_to_store", "g", "", "Path where you want store the file from one drive locally, example: -g /tmp/file.txt")
	GetfileOnedriveCmd.MarkFlagRequired("path_to_store")
}
