package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/bitfield/script"
	"github.com/charmbracelet/huh"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/sourcegraph/conc/pool"
)

func main() {
	ctx := context.Background()

	branches, err := openAndRefreshGit(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open/refresh repo: %v\n", err)
		os.Exit(1)
	}

	var createNewBranch bool

	form := huh.NewConfirm().
		Title("Create a new branch?").
		Affirmative("Yes, create").
		Negative("No, use existing").
		Value(&createNewBranch)

	if err := form.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "running branch confirmation form: %v\n", err)
		os.Exit(1)
	}

	var (
		branchInput huh.Field
		branchName  string
	)

	if createNewBranch {
		branchInput = huh.NewInput().
			Title("Name of new branch?").
			Validate(validateNewBranch).
			Value(&branchName)
	} else {
		branchInput = huh.NewSelect[string]().
			Title("Pick an existing branch.").
			Options(huh.NewOptions(branches...)...).
			Value(&branchName)
	}

	if err := branchInput.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "running branch input form: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Branch name: '%s'\n", branchName)

	fmt.Println("Done.")
}

func validateNewBranch(branchCandidate string) error {
	if !strings.HasPrefix(branchCandidate, "DEVPL-") && !strings.HasPrefix(branchCandidate, "NOJIRA/") {
		return fmt.Errorf("branch name '%s' should start with either 'DEVPL-' or 'NOJIRA/", branchCandidate)
	}

	return nil
}

// openAndRefreshGit returns a slice of short branch names, assuming nothing goes wrong.
func openAndRefreshGit(ctx context.Context) ([]string, error) {
	// Look for a git repo at or above the current working directory.
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("getting working directory: %w", err)
	}

	repo, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, fmt.Errorf("cannot open any directory at or above '%s' as a git repo: %w", wd, err)
	}

	// Pending implementation of: https://github.com/go-git/go-git/issues/74
	gitTopCmd := "git rev-parse --show-toplevel"

	gitTop, err := script.Exec(gitTopCmd).String()
	if err != nil {
		return nil, fmt.Errorf("executing '%s': %w", gitTopCmd, err)
	}

	gitTop = strings.TrimSpace(gitTop)

	// Update and prune all remotes using a worker pool.
	remotes, err := repo.Remotes()
	if err != nil {
		return nil, fmt.Errorf("reading remotes of repo at '%s': %w", gitTop, err)
	}

	fetchPool := pool.NewWithResults[string]().WithContext(ctx)

	for _, remote := range remotes {
		fetchPool.Go(func(ctx context.Context) (string, error) {
			remoteName := remote.Config().Name

			return remoteName, remote.FetchContext(ctx, &git.FetchOptions{
				RemoteName: remoteName,
				Prune:      true,
			})
		})
	}

	updatedRemotes, err := fetchPool.Wait()
	if err != nil {
		return nil, fmt.Errorf("updating/pruning remotes of repo at '%s': %w", gitTop, err)
	}

	fmt.Printf("Updated/pruned remotes: %s\n", strings.Join(updatedRemotes, ", "))

	refIter, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("getting branch references from repo at '%s': %w", gitTop, err)
	}

	branches := make([]string, 0)

	if err := refIter.ForEach(func(pRef *plumbing.Reference) error {
		branches = append(branches, pRef.Name().Short())
		return nil
	}); err != nil {
		return nil, fmt.Errorf("iterating over branch references in repo at '%s': %w", gitTop, err)
	}

	return branches, nil
}
