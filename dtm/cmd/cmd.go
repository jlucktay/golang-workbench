package cmd

import (
	"context"

	"go.jlucktay.dev/golang-workbench/dtm/docker"
)

func Execute(ctx context.Context, creator func() (docker.Client, error)) error {
	cli, err := creator()
	if err != nil {
		return err
	}

	_, err = cli.ListImages(ctx)
	if err != nil {
		return err
	}

	return cli.DeleteImage(ctx, "img string")
}
