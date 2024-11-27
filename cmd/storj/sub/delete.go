package sub

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/be"
	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
)

var (
	deleteSecret string
	deletePath   string
	deleteBucket string
)

// DeleteStorjCmd represents the delete command
var DeleteStorjCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete file from Storj bucket",
	Long:  "Description:\n  Delete a file from a Storj bucket based on the provided bucket name and path. The secret can be set as an environment variable named STORJ_SECRET for added security.\n\nExample:\n  sinaloa storj delete -s \"STORJ_SECRET\" -b \"storj_bucket_name\" -p \"folder1/version.json\"",
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LoadConfig()

		if deleteSecret == "" {
			deleteSecret = helpers.AppConfig.STORJ_SECRET
		}

		if deleteSecret == "" {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/delete' inside method 'DeleteStorjCmd': Required flag not set: secret",
				"500",
				nil,
			)
			return
		}

		storjService, err := be.NewStorjService(deleteSecret)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/delete' inside method 'DeleteStorjCmd' calling the function be 'NewStorjService(deleteSecret)'",
				"500",
				err,
			)
			return
		}

		err = storjService.DeleteFile(deleteBucket, deletePath)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/delete' inside method 'DeleteStorjCmd' calling the function be 'storjService.DeleteFile(deleteBucket, deletePath)'",
				"500",
				err,
			)
			return
		}

		successResponse := response.NewResponse(true, "200", "File deleted successfully from Storj", struct {
			Status string `json:"status"`
		}{
			Status: "ok",
		})
		successJsonResponse, err := json.MarshalIndent(successResponse, "", "  ")
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/delete' inside method 'DeleteStorjCmd' while marshaling JSON",
				"500",
				err,
			)
			return
		}

		fmt.Println(string(successJsonResponse))
	},
}

func init() {
	deleteSecretDefault := ""
	if helpers.AppConfig.STORJ_SECRET != "" {
		deleteSecretDefault = helpers.AppConfig.STORJ_SECRET
	}
	DeleteStorjCmd.Flags().StringVarP(&deleteSecret, "secret", "s", deleteSecretDefault, "Storj secret to connect with the bucket")
	DeleteStorjCmd.Flags().StringVarP(&deletePath, "path", "p", "", "Path of the file to delete on Storj bucket")
	DeleteStorjCmd.MarkFlagRequired("path")
	DeleteStorjCmd.Flags().StringVarP(&deleteBucket, "bucket", "b", "", "Name of the bucket")
	DeleteStorjCmd.MarkFlagRequired("bucket")
}
