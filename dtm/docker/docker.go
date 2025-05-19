package docker

import "context"

type Client interface {
	ListImages(ctx context.Context) ([]string, error)
	DeleteImage(ctx context.Context, img string) error
}
