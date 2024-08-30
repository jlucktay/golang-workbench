package main

import (
	"cmp"
	"context"
	"fmt"
	"maps"
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
	localImages := make([]localImage, 0)

	for index := range images {
		for rtIdx := range images[index].RepoTags {
			xrt := strings.Split(images[index].RepoTags[rtIdx], ":")

			newLI := localImage{
				id:   images[index].ID,
				repo: xrt[0],
				tag:  xrt[1],
			}

			localImages = append(localImages, newLI)
		}
	}

	// Sort parsed image data by repo and, if repo is equal, by tag (as semver).
	slices.SortStableFunc(localImages, func(a, b localImage) int {
		if n := cmp.Compare(a.repo, b.repo); n != 0 {
			return n
		} else if strings.EqualFold(a.tag, "latest") {
			return 1
		} else if strings.EqualFold(b.tag, "latest") {
			return -1
		}

		// If repo is equal, try parsing tag as semver and ordering that way.
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

	// Do a pretty print of the parsed image data.
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 5, 0, 1, ' ', 0)

	for index := range localImages {
		trimmedDigest := strings.TrimPrefix(localImages[index].id, "sha256:")

		fmt.Fprintf(w, "%s\t%s\t%s\n", localImages[index].repo, localImages[index].tag, trimmedDigest[:12])
	}

	w.Flush()
	fmt.Println()

	// Figure out which repos have more than one tag.
	tagOccurrences := make(map[string]int)

	for index := range localImages {
		tagOccurrences[localImages[index].repo]++
	}

	maps.DeleteFunc(tagOccurrences, func(_ string, v int) bool {
		return v <= 1
	})

	duplicateTags := make([]string, 0)

	for to := range tagOccurrences {
		duplicateTags = append(duplicateTags, to)
	}

	slices.Sort(duplicateTags)

	fmt.Println("Repos with more than one tag:")

	for _, dupe := range duplicateTags {
		fmt.Printf("  - %dx %s\n", tagOccurrences[dupe], dupe)
	}
}