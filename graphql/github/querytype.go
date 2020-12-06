package main

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
		} `graphql:"repositories(affiliations: OWNER, after: $endCursor, first: $perPage, isFork: $isFork, orderBy: {field: CREATED_AT, direction: ASC})"`
	} `graphql:"repositoryOwner(login: $login)"`
}
