package sub

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/controller"
)

var (
	repo         string
	itemsForPage = 100 // Default number of items per page
)

var GetImagesDockerCmd = &cobra.Command{
	Use:   "getImages",
	Short: "Short description of getImages",
	Long:  "Long description of getImages",
	Run: func(cmd *cobra.Command, args []string) {
		result, _ := controller.GetImages(repo, strconv.Itoa(itemsForPage))
		fmt.Println(string(result))
	},
}

func init() {
	GetImagesDockerCmd.Flags().IntVarP(&itemsForPage, "items", "i", itemsForPage, "Number of items per page")
	GetImagesDockerCmd.Flags().StringVarP(&repo, "repo", "r", "", "Docker repository to get images list")
	if err := GetImagesDockerCmd.MarkFlagRequired("repo"); err != nil {
		fmt.Println(err)
	}
}
