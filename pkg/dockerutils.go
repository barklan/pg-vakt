package pkg



import (
	"context"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func CheckIfNamedContainerIsRunning(nameInQuestion string) (bool, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false, err
	}

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return false, err
	}

	for _, container := range containers {
		log.Println("Container: ", container.Names)
		for _, name := range container.Names {
			if name[1:] == nameInQuestion {
				return true, nil
			}
		}
	}

	return false, nil
}
