package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/Masterminds/semver/v3"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type localImage struct {
	id   string
	repo string
	tag  string
}

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

	// Parse API response and gather up desired details on local images.
	mLocalImages := make(map[string][]localImage)
	imageRepos := make([]string, 0)

	for index := range images {
		for _, repoTag := range images[index].RepoTags {
			lastColon := strings.LastIndex(repoTag, ":")

			newLI := localImage{
				id:   images[index].ID,
				repo: repoTag[:lastColon],
				tag:  repoTag[lastColon+1:],
			}

			if _, inMap := mLocalImages[newLI.repo]; !inMap {
				mLocalImages[newLI.repo] = make([]localImage, 0)
				imageRepos = append(imageRepos, newLI.repo)
			}

			mLocalImages[newLI.repo] = append(mLocalImages[newLI.repo], newLI)
		}
	}

	slices.Sort(imageRepos)

	// Do a pretty print of the parsed image data.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 3, ' ', 0)

	for _, repo := range imageRepos {
		// Sort repo's images by tag (as semver).
		slices.SortStableFunc(mLocalImages[repo], func(a, b localImage) int {
			// Check if either tag is 'latest'.
			if strings.EqualFold(a.tag, "latest") {
				return 1
			} else if strings.EqualFold(b.tag, "latest") {
				return -1
			}

			// If neither tag is 'latest', try parsing tag as semver and ordering that way.
			verA, err := semver.NewVersion(a.tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse semver from '%s': %v\n", a.tag, err)
				return -1
			}

			verB, err := semver.NewVersion(b.tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse semver from '%s': %v\n", b.tag, err)
				return 1
			}

			return verA.Compare(verB)
		})

		for _, li := range mLocalImages[repo] {
			trimmedDigest := strings.TrimPrefix(li.id, "sha256:")
			fmt.Fprintf(w, "%s\t%s\t%s\n", repo, li.tag, trimmedDigest[:12])
		}
	}

	w.Flush()
}
