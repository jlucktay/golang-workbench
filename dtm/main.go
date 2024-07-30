package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func main() {
	ctx := context.Background()

	// Set up an API client to talk to the container engine (CE) server.
	cli, err := client.NewClientWithOpts(
		client.WithAPIVersionNegotiation(),
		client.WithHostFromEnv(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "creating new client: %v", err)
		return
	}
	defer cli.Close()

	// Request a list of images from the CE server.
	images, err := cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "listing images: %v", err)
		return
	}

	// 🚧 🚧 🚧
	for index := range images {
		fmt.Printf("%d: %#v\n", index, images[index])
	}
}
