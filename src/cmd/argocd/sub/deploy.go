package sub

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    json string
)

var DeployArgocdCmd = &cobra.Command{
    Use:   "deploy",
    Short: "Short description of deploy",
    Long:  "Long description of deploy",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Executing deploy in argocd")
    },
}

func init() {
    DeployArgocdCmd.Flags().StringVarP(&json, "json", "j", "", "Json to pass for the deploy with argocd")
}

