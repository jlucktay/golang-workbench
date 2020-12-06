package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func run(client *githubv4.Client, query *queryOwnedRepos, vars map[string]interface{}) ([]string, error) {
	var ownedRepos []string

	vars["endCursor"] = (*githubv4.String)(nil) // Null the 'after' argument to get first page.

	for {
		fmt.Printf("Querying with variables: %v... ", vars)

		if err := client.Query(context.TODO(), query, vars); err != nil {
			return nil, fmt.Errorf("couldn't run query: %w\n", err)
		}

		fmt.Printf("returned OK.\n")

		hasNextPage, endCursor := process(*query, &ownedRepos)
		if !hasNextPage {
			fmt.Println()
			break
		}

		vars["endCursor"] = githubv4.String(endCursor)
	}

	return ownedRepos, nil
}
