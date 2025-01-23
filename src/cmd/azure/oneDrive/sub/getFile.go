package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"
)

var (
	file          string
	path_to_store string
)

var GetfileOnedriveCmd = &cobra.Command{
	Use:   "get-file",
	Short: "Get a file from onedrive",
	Long:  "Get a file from onedrive",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the controller's GetFileList function
		result, _ := controller.GetFile(file, path_to_store)
		// Print the response
		fmt.Println(string(result))
	},
}

func init() {
	GetfileOnedriveCmd.Flags().StringVarP(&file, "file", "f", "", "file to get from onedrive, example: -f secrets/file.txt")
	GetfileOnedriveCmd.MarkFlagRequired("file")
	GetfileOnedriveCmd.Flags().StringVarP(&path_to_store, "path_to_store", "g", "", "path where you want store the file from onedrive locally, example: -g /tmp/file.txt")
	GetfileOnedriveCmd.MarkFlagRequired("path_to_store")
}
