package sub

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/errors"
	"github.com/eltiocaballoloco/sinaloa-cli/models/messages/response"
	"github.com/spf13/cobra"
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
		if statusCode, err := ping(urlPath); err != nil {
			errorMessage := fmt.Sprintf("Ping error: %s", err.Error())
			errorResponse := errors.NewErrorResponse(false, "500", errorMessage)
			errorJsonResponse, jsonErr := json.MarshalIndent(errorResponse, "", "  ")
			if jsonErr != nil {
				fmt.Println("Error marshaling JSON:", jsonErr)
				return
			}
			fmt.Println(string(errorJsonResponse))
			return
		} else {
			successData := struct {
				StatusCode int `json:"status_code"`
			}{
				StatusCode: statusCode,
			}
			successResponse := response.NewResponse(true, "200", "Ping successful", successData)
			jsonResponse, jsonErr := json.MarshalIndent(successResponse, "", "  ")
			if jsonErr != nil {
				fmt.Println("Error marshaling JSON:", jsonErr)
				return
			}
			fmt.Println(string(jsonResponse))
		}
	},
}

func ping(target string) (int, error) {
	var url string

	// Check if the target is an IP address
	if ip := net.ParseIP(target); ip != nil {
		// It's an IP address, so use it directly for pinging
		url = "http://" + target
	} else {
		// It's not an IP address, assume it's a domain name
		url = "http://" + target
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func init() {
	PingCmd.Flags().StringVarP(&urlPath, "url", "u", "", "The url to ping")
	if err := PingCmd.MarkFlagRequired("url"); err != nil {
		fmt.Println(err)
	}
}
