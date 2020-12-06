package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		"login":     githubv4.String(githubLogin),
		"perPage":   githubv4.Int(perPage),
		"endCursor": (*githubv4.String)(nil), // Null the 'after' argument to get first page.
	}

	var (
		query      queryOwnedRepos
		ownedRepos []string
	)

	for {
		fmt.Printf("Querying with variables: '%v'...", queryVariables)

		if err := client.Query(context.TODO(), &query, queryVariables); err != nil {
			fmt.Fprintf(os.Stderr, "couldn't run query: %v\n", err)
			return
		}

		fmt.Printf(" returned OK.\n")

		hasNextPage, endCursor := process(query, &ownedRepos)
		if !hasNextPage {
			break
		}

		queryVariables["endCursor"] = githubv4.String(endCursor)
	}

	fmt.Printf("\n%d owned repo(s):\n%s\n", len(ownedRepos), strings.Join(ownedRepos, "\n"))
}

type queryOwnedRepos *struct {
	RepositoryOwner *struct {
		Repositories *struct {
			TotalCount int
			PageInfo   *struct {
				HasNextPage bool
				EndCursor   string
			}
			Edges []*struct {
				Node *struct {
					Name string
				}
			}
		} `graphql:"repositories(affiliations: OWNER, after: $endCursor, first: $perPage, orderBy: {field: CREATED_AT, direction: ASC})"`
	} `graphql:"repositoryOwner(login: $login)"`
}

func process(qor queryOwnedRepos, ownedRepos *[]string) (hasNextPage bool, endCursor string) {
	if qor.RepositoryOwner != nil {
		if qor.RepositoryOwner.Repositories != nil {
			if qor.RepositoryOwner.Repositories.PageInfo != nil {
				hasNextPage = qor.RepositoryOwner.Repositories.PageInfo.HasNextPage
				endCursor = qor.RepositoryOwner.Repositories.PageInfo.EndCursor
			}

			if qor.RepositoryOwner.Repositories.Edges != nil {
				for i := range qor.RepositoryOwner.Repositories.Edges {
					if qor.RepositoryOwner.Repositories.Edges[i].Node != nil {
						*ownedRepos = append(*ownedRepos, qor.RepositoryOwner.Repositories.Edges[i].Node.Name)
					}
				}
			}
		}
	}

	return
}
