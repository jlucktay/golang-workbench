package main

func filterForGoRepos(input map[string]int) bool {
	if _, ok := input["Go"]; ok {
		for _, lineCount := range input {
			if lineCount > input["Go"] {
				return false
			}
		}

		return true
	}

	return false
}
