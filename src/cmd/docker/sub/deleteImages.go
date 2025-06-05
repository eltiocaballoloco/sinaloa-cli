package sub

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/controller"
)

var (
	repoD         string
	itemsToTake   int
	itemsForPageD int    = 100    // Default number of items per page
	dryRunStr     string = "true" // Default to dry run
)

var DeleteImagesDockerCmd = &cobra.Command{
	Use:   "delete-images",
	Short: "Delete Docker images from a repository",
	Long:  "Delete Docker images from a specified repository on Docker Hub. You can specify the number of items per page.",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, err := strconv.ParseBool(dryRunStr)
		if err != nil {
			fmt.Println("[Error] Invalid value for dry-run, must be true or false")
			return
		}
		result, _ := controller.GetImages(repoD, strconv.Itoa(itemsForPageD), strconv.Itoa(itemsToTake), "delete", dryRun)
		fmt.Println(string(result))
	},
}

func init() {
	DeleteImagesDockerCmd.Flags().IntVarP(&itemsForPageD, "items", "i", itemsForPageD, "Number of items per page")
	DeleteImagesDockerCmd.Flags().StringVarP(&repoD, "repo", "r", "", "Docker repository to get images list")
	if err := DeleteImagesDockerCmd.MarkFlagRequired("repo"); err != nil {
		fmt.Println(err)
	}
	DeleteImagesDockerCmd.Flags().IntVarP(&itemsToTake, "items-to-take", "t", 0, "Number of the first X tags to take from the repository.")
	DeleteImagesDockerCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed("items-to-take") || itemsToTake <= 0 {
			return fmt.Errorf("parameter '--items-to-take' is required and must be greater than 0")
		}
		return nil
	}
	DeleteImagesDockerCmd.Flags().StringVarP(&dryRunStr, "dry-run", "d", dryRunStr, "Dry run, true or false")
}
