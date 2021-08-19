package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func prettyPrint(input []string) {
	fmt.Printf("%d owned repo(s):\n", len(input))

	// get terminal width
	tw, _, errTGS := term.GetSize(int(os.Stdout.Fd()))
	if errTGS != nil {
		fmt.Fprintf(os.Stderr, "couldn't get terminal size: %v\n", errTGS)

		return
	}

	// get longest repo name
	longestRepoName := 0
	for i := range input {
		if len(input[i]) > longestRepoName {
			longestRepoName = len(input[i])
		}
	}

	// do math to divide lines evenly across width
	longestRepoName++ // add a single padding space
	reposPerLine := tw / longestRepoName

	// space out repo names in columns and pretty print
	for i := 0; i < len(input); i += reposPerLine {
		for j := 0; j < reposPerLine && i+j < len(input); j++ {
			fmt.Printf("%-[1]*[2]s", longestRepoName, input[i+j])
		}

		fmt.Println()
	}

	fmt.Println()
}
