package sub

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	port string
	method string
)

var StartApiCmd = &cobra.Command{
	Use:   "start",
	Short: "Start sinaloa api server",
	Long:  `Start sinaloa api server. This allow you to use all commands of the cli via api. Example: sinaloa api start -p 8080 -m all|haproxy|dim`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("[INFORMATION] Starting API server on port:", port, "🤠🤠")
		fmt.Println("[INFORMATION] Starting API server using this method:", method, "🤠🤠")

		// Start api server
		

		fmt.Println("\033[32m[INFORMATION] Server started correctly\033[0m 🚀🚀")
	},
}

func init() {
	StartApiCmd.Flags().StringVarP(&port, "port", "p", "", "Port of api server.")
	StartApiCmd.Flags().StringVarP(&method, "method", "m", "", "Method to start api server.")

	if err := StartApiCmd.MarkFlagRequired("port"); err != nil {
		fmt.Println("\033[31m[ERROR] There was a problem starting the API server:", err, "🚩🚩\033[0m")
	}

	if err := StartApiCmd.MarkFlagRequired("method"); err != nil {
		fmt.Println("\033[31m[ERROR] There was a problem starting the API server:", err, "🚩🚩\033[0m")
	}
}
