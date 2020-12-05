package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

func main() {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	httpClient := oauth2.NewClient(context.Background(), src)

	client := githubv4.NewClient(httpClient)

	if err := client.Query(context.Background(), &query, nil); err != nil {
		fmt.Fprintf(os.Stderr, "couldn't run query: %v", err)
	}

	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)
}

var query struct {
	Viewer struct {
		Login     githubv4.String
		CreatedAt githubv4.String
	}
}
