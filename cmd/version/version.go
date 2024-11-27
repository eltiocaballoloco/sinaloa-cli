package version

import (
	"encoding/json"
	"fmt"

	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get the version of sinaloa-cli",
	Long:  "Get the version of sinaloa-cli. Example: sinaloa version",
	Run: func(cmd *cobra.Command, args []string) {
		// Load the config
		helpers.LoadConfig()

		// Create the response
		versionResponse := response.NewResponse(true, "200", "", struct {
			Version string `json:"version"`
		}{
			Version: helpers.AppConfig.VERSION,
		})

		// Marshal the response into JSON with indentation
		jsonResponse, err := json.MarshalIndent(versionResponse, "", "  ")
		if err != nil {
			errorMessage := fmt.Sprintf("An error occurred: %s", err.Error())
			errorResponse := errors.NewErrorResponse(false, "500", errorMessage)

			errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
			if jsonErr != nil {
				fmt.Println("Error marshaling the error message:", jsonErr)
				return
			}

			fmt.Println(string(errorJsonResponse))
			return
		}

		fmt.Println(string(jsonResponse))
	},
}

func init() {}
