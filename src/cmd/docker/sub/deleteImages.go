package sub

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/controller"
)

var (
	repoD         string
	dryRun        bool = true // Default to dry run
	itemsForPageD int  = 100  // Default number of items per page
)

var DeleteImagesDockerCmd = &cobra.Command{
	Use:   "delete-images",
	Short: "Delete Docker images from a repository",
	Long:  "Delete Docker images from a specified repository on Docker Hub. You can specify the number of items per page.",
	Run: func(cmd *cobra.Command, args []string) {
		result, _ := controller.GetImages(repo, strconv.Itoa(itemsForPage), "delete", dryRun)
		fmt.Println(string(result))
	},
}

func init() {
	DeleteImagesDockerCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", dryRun, "Dry run, if true return only the tags to delete without actually deleting them. False will delete the tags")
	DeleteImagesDockerCmd.Flags().IntVarP(&itemsForPageD, "items", "i", itemsForPageD, "Number of items per page")
	DeleteImagesDockerCmd.Flags().StringVarP(&repoD, "repo", "r", "", "Docker repository to get images list")
	if err := DeleteImagesDockerCmd.MarkFlagRequired("repo"); err != nil {
		fmt.Println(err)
	}
}
