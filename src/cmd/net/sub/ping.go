package sub

import (
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/net/controller"
)

var urlPath string
var client = &http.Client{
	Timeout: 2 * time.Second,
}

var PingCmd = &cobra.Command{
	Use:   "ping",
	Short: "This command is used to ping a url or an ip address",
	Long:  `This command is used to ping a url or an ip address. Return 200 if ping it is ok otherwise error. Example: sinaloa net ping -u google.com`,
	Run: func(cmd *cobra.Command, args []string) {
		// Call the Ping controller
		result, _ := controller.Ping(urlPath)
		// Print the success response
		fmt.Println(string(result))
	},
}

func init() {
	PingCmd.Flags().StringVarP(&urlPath, "url", "u", "", "the url to ping, eg. google.com")
	if err := PingCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
}
