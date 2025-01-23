package sub

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/azure/oneDrive/controller"
)

var (
	file_path_to_upload string
	upload_path         string
)

var UploadFileOnedriveCmd = &cobra.Command{
	Use:   "upload-file",
	Short: "Upload a file to onedrive",
	Long:  "Upload a file to onedrive",
	Run: func(cmd *cobra.Command, args []string) {
		// Call the controller's UploadFile function
		result, _ := controller.UploadFile(file_path_to_upload, upload_path)
		// Print the response
		fmt.Println(string(result))
	},
}

func init() {
	UploadFileOnedriveCmd.Flags().StringVarP(&file_path_to_upload, "file_path_to_upload", "f", "", "file to upload to onedrive, example: -f /local/path/file.txt")
	UploadFileOnedriveCmd.MarkFlagRequired("file_path_to_upload")
	UploadFileOnedriveCmd.Flags().StringVarP(&upload_path, "upload_path", "g", "", "path where you want store the file in onedrive, example: -g config/file.txt")
	UploadFileOnedriveCmd.MarkFlagRequired("upload_path")
}
