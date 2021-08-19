package main

type queryOwnedRepos *struct {
	RepositoryOwner *struct {
		Repositories *struct {
			PageInfo *struct {
				EndCursor   string
				HasNextPage bool
			}
			Edges []*struct {
				Node *struct {
					Name string
				}
			}
			TotalCount int
		} `graphql:"repositories(affiliations: OWNER, after: $endCursor, first: $perPage, isFork: $isFork, orderBy: {field: CREATED_AT, direction: ASC})"`
	} `graphql:"repositoryOwner(login: $login)"`
}
