package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/go-git/go-git/v5"
	"github.com/sourcegraph/conc/pool"
)

func main() {
	// Look for a git repo at or above the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "getting working directory: %v\n", err)
		os.Exit(1)
	}

	repo, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open any directory at or above '%s' as a git repo: %v\n", wd, err)
		os.Exit(1)
	}

	// Pending implementation of: https://github.com/go-git/go-git/issues/74
	gitTopCmd := "git rev-parse --show-toplevel"

	gitTop, err := script.Exec(gitTopCmd).String()
	if err != nil {
		fmt.Fprintf(os.Stderr, "executing '%s': %v\n", gitTopCmd, err)
		os.Exit(1)
	}

	gitTop = strings.TrimSpace(gitTop)

	// Update and prune all remotes using a worker pool.
	remotes, err := repo.Remotes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "reading remotes of repo at '%s': %v\n", gitTop, err)
		os.Exit(1)
	}

	ctx := context.Background()
	fetchPool := pool.New().WithContext(ctx)

	for _, remote := range remotes {
		fetchPool.Go(func(ctx context.Context) error {
			remoteName := remote.Config().Name

			fmt.Printf("fetching remote '%s'...\n", remoteName)

			return remote.FetchContext(ctx, &git.FetchOptions{
				RemoteName: remoteName,
				Prune:      true,
			})
		})
	}

	if err := fetchPool.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "fetching remotes of repo at '%s': %v\n", gitTop, err)
		os.Exit(1)
	}
}
