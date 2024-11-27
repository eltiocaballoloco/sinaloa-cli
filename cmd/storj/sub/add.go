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
	secret      string
	path        string
	file        string
	bucket_name string
	fileContent string
)

var AddStorjCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file to storj bucket",
	Long:  `Description: Add file to storj bucket. The secret can be an env os variable called STORJ_SECRET, this to add more security without passing by args the secret. Example: sinaloa storj add -s "STORJ_SECRET" -b "storj_bucket_name" -p "folder1/version.json" -f "/Users/user1/version.json" -c "content"`,
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LoadConfig()

		if secret == "" && helpers.AppConfig.STORJ_SECRET == "" {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/add' inside method 'AddStorjCmd': Required flag not set: secret",
				"500",
				nil,
			)
			return
		}

		if helpers.AppConfig.STORJ_SECRET != "" && secret == "" {
			secret = helpers.AppConfig.STORJ_SECRET
		}

		storjService, err := be.NewStorjService(secret)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/add' inside method 'AddStorjCmd' calling the function be 'be.NewStorjService(secret)'",
				"401",
				err,
			)
			return
		}

		if fileContent != "" {
			err = storjService.AddFileBytes(bucket_name, path, []byte(fileContent))
		} else {
			err = storjService.AddFile(bucket_name, path, file)
		}
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/add' inside method 'AddStorjCmd' calling the function be 'uplink.storjService.AddFile'",
				"500",
				err,
			)
			return
		}

		successResponse := response.NewResponse(true, "200", "File uploaded successfully to Storj", struct {
			Status string `json:"status"`
		}{Status: "ok"})
		successJsonResponse, err := json.MarshalIndent(successResponse, "", "  ")
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/add' inside method 'AddStorjCmd' while marshaling JSON",
				"500",
				err,
			)
			return
		}

		fmt.Println(string(successJsonResponse))
	},
}

func init() {
	secretDefault := ""
	if helpers.AppConfig.STORJ_SECRET != "" {
		secretDefault = helpers.AppConfig.STORJ_SECRET
	}
	AddStorjCmd.Flags().StringVarP(&secret, "secret", "s", secretDefault, "Storj secret to connect with the bucket")
	AddStorjCmd.Flags().StringVarP(&path, "path", "p", "", "Path where you want store the file on storj bucket")
	AddStorjCmd.Flags().StringVarP(&file, "file", "f", "", "Path of the file where is located")
	AddStorjCmd.Flags().StringVarP(&bucket_name, "bucket", "b", "", "Name of your bucket")
	AddStorjCmd.Flags().StringVarP(&fileContent, "content", "c", "", "Direct content to be stored in the file")

	AddStorjCmd.MarkFlagRequired("path")
	AddStorjCmd.MarkFlagRequired("bucket")
}
