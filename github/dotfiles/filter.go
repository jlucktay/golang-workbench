package main

func filterForGoRepos(input map[string]int) bool {
	i, ok := input["Go"]
	if !ok {
		return false
	}

	for _, lineCount := range input {
		if lineCount > i {
			return false
		}
	}

	return true
}
