package main

import (
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func TestActivityNotifications(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	inputBytes, err := os.ReadFile("input02.txt")
	is.NoErr(err)

	inputLines := strings.Split(string(inputBytes), "\n")
	is.True(len(inputLines) == 3)

	firstLine := strings.Fields(inputLines[0])
	is.True(len(firstLine) == 2)

	numDaysData := firstLine[0]
	_, err = strconv.Atoi(numDaysData)
	is.NoErr(err)

	numTrailingDays := firstLine[1]
	ntd, err := strconv.Atoi(numTrailingDays)
	is.NoErr(err)

	expenditure := make([]int, 0)

	secondLine := strings.Fields(inputLines[1])

	for i := range secondLine {
		expenditureItem, err := strconv.Atoi(secondLine[i])
		is.NoErr(err)

		expenditure = append(expenditure, expenditureItem)
	}

	expectedOutputBytes, err := os.ReadFile("output02.txt")
	is.NoErr(err)

	expectedOutput, err := strconv.Atoi(strings.TrimSpace(string(expectedOutputBytes)))
	is.NoErr(err)

	actualOutput := activityNotifications(expenditure, ntd)
	is.Equal(expectedOutput, actualOutput)
}
