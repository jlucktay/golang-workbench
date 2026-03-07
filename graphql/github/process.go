package main

func process(qor queryOwnedRepos, ownedRepos *[]string) (bool, string) {
	if qor.RepositoryOwner == nil || qor.RepositoryOwner.Repositories == nil || qor.RepositoryOwner.Repositories.PageInfo == nil {
		return false, ""
	}

	hasNextPage := qor.RepositoryOwner.Repositories.PageInfo.HasNextPage
	endCursor := qor.RepositoryOwner.Repositories.PageInfo.EndCursor

	if qor.RepositoryOwner.Repositories.Edges == nil {
		return hasNextPage, endCursor
	}

	for i := range qor.RepositoryOwner.Repositories.Edges {
		if qor.RepositoryOwner.Repositories.Edges[i].Node != nil {
			*ownedRepos = append(*ownedRepos, qor.RepositoryOwner.Repositories.Edges[i].Node.Name)
		}
	}

	return hasNextPage, endCursor
}
