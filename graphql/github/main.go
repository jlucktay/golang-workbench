package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const (
	envKey      = "GITHUB_TOKEN"
	githubLogin = "jlucktay"
	perPage     = 100

	timeout5s = 5 * time.Second
)

func main() {
	// Set up GitHub GraphQL API v4 client
	token, tokenSet := os.LookupEnv(envKey)
	if !tokenSet {
		fmt.Fprintf(os.Stderr, "token not set in environment: %s\n", envKey)

		return
	}

	oaToken := &oauth2.Token{AccessToken: token} //nolint:exhaustivestruct
	src := oauth2.StaticTokenSource(oaToken)

	ctx, cancel := context.WithTimeout(context.Background(), timeout5s)
	defer cancel()

	httpClient := oauth2.NewClient(ctx, src)
	httpClient.Timeout = timeout5s
	client := githubv4.NewClient(httpClient)

	// Set up queries to send to GraphQL and hold results, and variables for each run
	var queryMine, queryForked queryOwnedRepos

	queryVariables := map[string]interface{}{
		"login":   githubv4.String(githubLogin),
		"perPage": githubv4.Int(perPage),
	}

	// Query for unforked repos
	queryVariables["isFork"] = githubv4.Boolean(false)

	myRepos, errRunMine := runQuery(client, &queryMine, queryVariables)
	if errRunMine != nil {
		fmt.Fprint(os.Stderr, errRunMine)

		return
	}

	prettyPrint(myRepos)

	// Query for forked repos
	queryVariables["isFork"] = githubv4.Boolean(true)

	forkedRepos, errRunForked := runQuery(client, &queryForked, queryVariables)
	if errRunForked != nil {
		fmt.Fprint(os.Stderr, errRunForked)

		return
	}

	prettyPrint(forkedRepos)
}
