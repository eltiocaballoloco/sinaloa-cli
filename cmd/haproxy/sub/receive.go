package sub

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    msg string
)

var ReceiveHaproxyCmd = &cobra.Command{
    Use:   "receive",
    Short: "This command is used to execute receive in haproxy.",
    Long:  "This command is used to execute receive in haproxy. Is used for example, to renew SSL certificates, if haproxy architecture implement multiple haproxy endpoints sharinga virtual ip. Example: sinaloa haproxy receive -m \"Message to receive\"",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Executing receive in haproxy")
    },
}

func init() {
    ReceiveHaproxyCmd.Flags().StringVarP(&msg, "msg", "m", "", "Message to receive")
    ReceiveHaproxyCmd.MarkFlagRequired("msg")
}

