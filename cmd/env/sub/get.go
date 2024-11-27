package sub

import (
	"errors"
	"fmt"
	"os"

	"github.com/eltiocaballoloco/sinaloa-cli/helpers"
	"github.com/spf13/cobra"
)

var (
	key string
)

var GetEnvCmd = &cobra.Command{
	Use:   "get",
	Short: "Get env variables",
	Long:  "Get env variables. Example: sinaloa env get -k \"KEY\"",
	Run: func(cmd *cobra.Command, args []string) {
		// Get the value of the env variable
		value := os.Getenv(key)
		// Check if the value is empty
		if value == "" {
			err := errors.New(fmt.Sprintf("Environment variable %s is not set", key))
			helpers.HandleResponseError(
				"An error occurred",
				"404",
				err,
			)
			return
		} else {
			helpers.HandleResponse(
				"Variable found",
				"200",
				map[string]string{
					"key":   key,
					"value": value,
				},
			)
			return
		}
	},
}

func init() {
	GetEnvCmd.Flags().StringVarP(&key, "key", "k", "", "Env variable name to get value")
	GetEnvCmd.MarkFlagRequired("key")
}
