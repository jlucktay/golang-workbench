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
	queryVariables := map[string]interface{}{"login": githubv4.String(githubLogin)}

	if err := client.Query(context.TODO(), &query, queryVariables); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run query: %v\n", err)
		return
	}

	if query.RepositoryOwner != nil {
		if query.RepositoryOwner.Repositories != nil {
			fmt.Printf("TotalCount: %d\n", query.RepositoryOwner.Repositories.TotalCount)

			if query.RepositoryOwner.Repositories.Edges != nil {
				for i := range query.RepositoryOwner.Repositories.Edges {
					if query.RepositoryOwner.Repositories.Edges[i].Node != nil {
						fmt.Printf("%d\t", i)
						fmt.Printf("CreatedAt: '%s'\t", query.RepositoryOwner.Repositories.Edges[i].Node.CreatedAt)
						fmt.Printf("Name: '%s'\n", query.RepositoryOwner.Repositories.Edges[i].Node.Name)
					}
				}
			}
		}
	}
}

var query *struct {
	RepositoryOwner *struct {
		Login        githubv4.String
		Repositories *struct {
			TotalCount githubv4.Int
			Edges []*struct {
				Node *struct {
					CreatedAt githubv4.String
					Name      githubv4.String
				}
			}
		} `graphql:"repositories(affiliations: OWNER, first: 10, orderBy: {field: CREATED_AT, direction: ASC})"`
	} `graphql:"repositoryOwner(login: $login)"`
}
