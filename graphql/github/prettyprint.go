package main

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

func prettyPrintTerminal(input []string, repoType string) {
	fmt.Printf("%d repo %s:\n", len(input), repoType)

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

type jsonOutput struct {
	Sources []string `json:"sources"`
	Forks   []string `json:"forks"`
}

var jsonBuffer jsonOutput

const (
	printSources = "sources"
	printForks   = "forks"
)

func prettyPrintJSON(input []string, repoType string) {
	switch repoType {
	case printSources:
		jsonBuffer.Sources = append(jsonBuffer.Sources, input...)
	case printForks:
		jsonBuffer.Forks = append(jsonBuffer.Forks, input...)
	}
}
