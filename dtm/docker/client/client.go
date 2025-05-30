package client

import (
	"context"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"

	"go.jlucktay.dev/golang-workbench/dtm/docker"
)

type Docker struct {
	client.APIClient
}

// New creates a new Docker client using sensible defaults.
func New() (docker.Client, error) {
	//  ...
	return &Docker{}, nil
}

func (d *Docker) ListImages(ctx context.Context) ([]string, error) {
	_, err := d.ImageList(ctx, image.ListOptions{})

	return nil, err
}

func (d *Docker) DeleteImage(ctx context.Context, img string) error {
	_, err := d.ImageRemove(ctx, img, image.RemoveOptions{})

	return err
}
