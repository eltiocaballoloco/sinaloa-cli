package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/azure/oneDrive/controller"
)

var (
	path string
)

var GetfileListOnedriveCmd = &cobra.Command{
	Use:   "get-file-list",
	Short: "Get a list of file and folders from one drive",
	Long:  "Get a list of file and folders from one drive",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the controller's GetFileList function
		fileList, err := controller.GetFileList(path)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return
		}

		// Print the list of files and folders
		fmt.Println("Files and Folders in OneDrive:")
		for _, item := range fileList {
			fmt.Println(item)
		}
	},
}

func init() {
	GetfileListOnedriveCmd.Flags().StringVarP(&path_to_store, "path", "g", "", "path where you want see the list of the files or folders, example: -g /docs or . to show the root")
	GetfileListOnedriveCmd.MarkFlagRequired("path")
}