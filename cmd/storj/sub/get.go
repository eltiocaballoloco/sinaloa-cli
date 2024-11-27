package sub

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/be"
	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
)

var (
	getSecret    string
	getPath      string
	getBucket    string
	getPathLocal string
)

// GetStorjCmd represents the get command for fetching a file from Storj
var GetStorjCmd = &cobra.Command{
	Use:   "get",
	Short: "Get file from Storj bucket",
	Long:  "Description:\n  Retrieve a file from a specified Storj bucket and path.\n  The secret can be an env os variable called STORJ_SECRET,\n  this to add more security without passing by args the secret.\n\nExample:\n - sinaloa storj get -s \"STORJ_SECRET\" -b \"storj_bucket_name\" -p \"folder1/version.json\"\n - sinaloa storj get -s \"STORJ_SECRET\" -b \"storj_bucket_name\" -p \"folder1/version.json\" -m \"/Users/user1/version.json\"",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the configuration
		helpers.LoadConfig()

		// Use the provided secret or fallback to the environment variable
		if getSecret == "" {
			getSecret = helpers.AppConfig.STORJ_SECRET
		}

		// Check if the secret is empty
		if getSecret == "" {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd': Required flag not set: secret",
				"500",
				nil,
			)
			return
		}

		// Initialize Storj service
		storjService, err := be.NewStorjService(getSecret)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' calling the function be 'uplink.NewStorjService(getSecret)'",
				"500",
				err,
			)
			return
		}

		// Fetch file from Storj
		data, err := storjService.GetFile(getBucket, getPath)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' calling the function be 'uplink.storjService.GetFile(getBucket, getPath)'",
				"500",
				err,
			)
			return
		}

		// Determine the file type
		fileType := http.DetectContentType(data)

		// Get file name
		fileName, err := be.ExtractFileName(getPath)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' calling the function be 'utils.ExtractFileName(getPath)'",
				"500",
				err,
			)
			return
		}

		// Get file extension
		fileExtension, err := be.ExtractFileExtension(getPath)
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' calling the function be 'utils.ExtractFileExtension(getPath)'",
				"500",
				err,
			)
			return
		}

		// Save the file locally if path-local is provided
		var fileLocallySavedAt string
		if getPathLocal != "" {
			err = be.SaveFileLocally(getPathLocal, data)
			if err != nil {
				helpers.HandleResponseError(
					"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' calling the function be 'utils.SaveFileLocally(getPathLocal, data)'",
					"500",
					err,
				)
				return
			}

			fileLocallySavedAt = getPathLocal
			fmt.Printf("File saved successfully at local path: %s\n", getPathLocal)
		}

		// Success response
		successMessage := "File fetched successfully from Storj"
		responseContent := struct {
			FileType           string `json:"file_type"`
			FileName           string `json:"file_name"`
			FileExtension      string `json:"file_extension"`
			FileLocallySavedAt string `json:"file_locally_saved_at,omitempty"`
			FileContents       string `json:"file_contents,omitempty"`
		}{
			FileType:           fileType,
			FileName:           fileName,
			FileExtension:      fileExtension,
			FileLocallySavedAt: fileLocallySavedAt,
		}

		if getPathLocal == "" {
			responseContent.FileContents = string(data)
		}

		successResponse := response.NewResponse(true, "200", successMessage, responseContent)
		successJsonResponse, err := json.MarshalIndent(successResponse, "", "  ")
		if err != nil {
			helpers.HandleResponseError(
				"An error occurred at 'github.com/eltiocaballoloco/sinaloa-cli/cmd/storj/sub/get' inside method 'GetStorjCmd' while marshaling JSON",
				"500",
				err,
			)
			return
		}

		fmt.Println(string(successJsonResponse))
	},
}

func init() {
	getSecretDefault := ""
	if helpers.AppConfig.STORJ_SECRET != "" {
		getSecretDefault = helpers.AppConfig.STORJ_SECRET
	}
	GetStorjCmd.Flags().StringVarP(&getSecret, "secret", "s", getSecretDefault, "Storj secret to connect with the bucket")
	GetStorjCmd.Flags().StringVarP(&getPath, "path", "p", "", "Path of the file to get from the storj bucket")
	GetStorjCmd.Flags().StringVarP(&getBucket, "bucket", "b", "", "Name of the storj bucket")
	GetStorjCmd.Flags().StringVarP(&getPathLocal, "path-local", "m", "", "Local path to save the file")

	GetStorjCmd.MarkFlagRequired("path")
	GetStorjCmd.MarkFlagRequired("bucket")
}
