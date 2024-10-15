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
	ctx := context.Background()

	if err := openAndRefreshGit(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "could not open/refresh repo: %v\n", err)
		os.Exit(1)
	}
}

func openAndRefreshGit(ctx context.Context) error {
	// Look for a git repo at or above the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	repo, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return fmt.Errorf("cannot open any directory at or above '%s' as a git repo: %w", wd, err)
	}

	// Pending implementation of: https://github.com/go-git/go-git/issues/74
	gitTopCmd := "git rev-parse --show-toplevel"

	gitTop, err := script.Exec(gitTopCmd).String()
	if err != nil {
		return fmt.Errorf("executing '%s': %w", gitTopCmd, err)
	}

	gitTop = strings.TrimSpace(gitTop)

	// Update and prune all remotes using a worker pool.
	remotes, err := repo.Remotes()
	if err != nil {
		return fmt.Errorf("reading remotes of repo at '%s': %w", gitTop, err)
	}

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
		return fmt.Errorf("fetching remotes of repo at '%s': %w", gitTop, err)
	}

	return nil
}
