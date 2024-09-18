package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"text/tabwriter"

	"github.com/Masterminds/semver/v3"
	"github.com/charmbracelet/huh"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type localImage struct {
	id   string
	repo string
	tag  string
}

var patchZero = regexp.MustCompile(`\.\d+\.0`)

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

	deletionCandidates := make([]localImage, 0)

	for _, repo := range imageRepos {
		if len(mLocalImages[repo]) <= 1 {
			continue
		}

		for index, li := range mLocalImages[repo] {
			if strings.EqualFold(li.tag, "latest") {
				continue
			}

			sv, err := semver.NewVersion(li.tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse semver from '%s': %v\n", li.tag, err)
				continue
			}

			if len(mLocalImages[repo]) <= index+1 {
				continue
			}

			nextLI := mLocalImages[repo][index+1]

			if strings.EqualFold(nextLI.tag, "latest") {
				continue
			}

			hv, err := semver.NewVersion(nextLI.tag)
			if err != nil {
				fmt.Fprintf(os.Stderr, "could not parse semver from '%s': %v\n", nextLI.tag, err)
				continue
			}

			if (hv.Major() > sv.Major() ||
				(hv.Minor() > sv.Minor() && sv.Minor() > 0) ||
				(hv.Patch() > sv.Patch() && (sv.Patch() > 0 || patchZero.Match([]byte(sv.Original()))))) &&
				len(hv.Prerelease()) == len(sv.Prerelease()) && len(hv.Metadata()) == len(sv.Metadata()) {

				deletionCandidates = append(deletionCandidates, li)
			}
		}
	}

	toBeDeleted := make([]localImage, 0)

	for _, dc := range deletionCandidates {
		var confirm bool

		title := fmt.Sprintf("Delete '%s:%s' image?", dc.repo, dc.tag)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title(title).
					Affirmative("Yes!").
					Negative("No.").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "running form: %v", err)
			continue
		}

		if confirm {
			toBeDeleted = append(toBeDeleted, dc)
		}
	}

	if len(toBeDeleted) > 0 {
		fmt.Println("Deleting images:")
	}

	for index, tbd := range toBeDeleted {
		deleteImageTag := fmt.Sprintf("%s:%s", tbd.repo, tbd.tag)

		fmt.Printf("[%d] %s...\n", index, deleteImageTag)

		removeResult, err := cli.ImageRemove(ctx, deleteImageTag, image.RemoveOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "removing image '%s': %v", deleteImageTag, err)
			continue
		}

		for _, dr := range removeResult {
			if len(dr.Deleted) > 0 {
				fmt.Printf("  Deleted '%s'.\n", dr.Deleted)
			}
			if len(dr.Untagged) > 0 {
				fmt.Printf("  Untagged '%s'.\n", dr.Untagged)
			}
		}

		fmt.Println("Done.")
	}
}
