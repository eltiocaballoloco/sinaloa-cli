package docker

import (
	"github.com/eltiocaballoloco/sinaloa-cli/src/cmd/docker/sub"
	"github.com/spf13/cobra"
)

var DockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker commands",
	Long:  "Docker commands to manage registries repos, images, containers and other docker resources",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	DockerCmd.AddCommand(sub.GetImagesDockerCmd)
	DockerCmd.AddCommand(sub.DeleteImagesDockerCmd)
}
