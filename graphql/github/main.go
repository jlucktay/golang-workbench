package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	envToken    = "GITHUB_TOKEN"
	githubLogin = "jlucktay"
	perPage     = 10
)

func main() {
	token, tokenSet := os.LookupEnv(envToken)
	if !tokenSet {
		fmt.Fprintf(os.Stderr, "token not set in environment: %s\n", envToken)
		return
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	httpClient := oauth2.NewClient(context.TODO(), src)
	client := githubv4.NewClient(httpClient)
	queryVariables := map[string]interface{}{
		"endCursor": (*githubv4.String)(nil), // Null the 'after' argument to get first page.
		"login":     githubv4.String(githubLogin),
		"perPage":   githubv4.Int(perPage),
	}

	var (
		query      queryOwnedRepos
		ownedRepos []string
	)

	for {
		fmt.Printf("Querying with variables: %v...", queryVariables)

		if err := client.Query(context.TODO(), &query, queryVariables); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't run query: %v\n", err)
			return
		}

		fmt.Printf(" returned OK.\n")

		hasNextPage, endCursor := process(query, &ownedRepos)
		if !hasNextPage {
			fmt.Println()
			break
		}

		queryVariables["endCursor"] = githubv4.String(endCursor)
	}

	prettyPrint(ownedRepos)
}
