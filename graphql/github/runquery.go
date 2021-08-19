package main

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
)

func runQuery(client *githubv4.Client, query *queryOwnedRepos, vars map[string]interface{}) ([]string, error) {
	var (
		ownedRepos []string
		endCursor  string

		hasNextPage = true
	)

	vars["endCursor"] = (*githubv4.String)(nil) // Null the 'after' argument to get first page.

	for hasNextPage {
		fmt.Printf("Querying with variables: %v... ", vars)

		if err := client.Query(context.TODO(), query, vars); err != nil {
			return nil, fmt.Errorf("couldn't run query: %w", err)
		}

		fmt.Println("returned OK.")

		hasNextPage, endCursor = process(*query, &ownedRepos)

		vars["endCursor"] = githubv4.String(endCursor)
	}

	fmt.Println()

	return ownedRepos, nil
}
