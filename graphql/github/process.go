package main

func process(qor queryOwnedRepos, ownedRepos *[]string) (hasNextPage bool, endCursor string) {
	if qor.RepositoryOwner == nil {
		return
	}

	if qor.RepositoryOwner.Repositories == nil {
		return
	}

	if qor.RepositoryOwner.Repositories.PageInfo == nil {
		return
	}

	hasNextPage = qor.RepositoryOwner.Repositories.PageInfo.HasNextPage
	endCursor = qor.RepositoryOwner.Repositories.PageInfo.EndCursor

	if qor.RepositoryOwner.Repositories.Edges == nil {
		return
	}

	for i := range qor.RepositoryOwner.Repositories.Edges {
		if qor.RepositoryOwner.Repositories.Edges[i].Node != nil {
			*ownedRepos = append(*ownedRepos, qor.RepositoryOwner.Repositories.Edges[i].Node.Name)
		}
	}

	return
}
