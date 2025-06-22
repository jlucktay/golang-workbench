package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
 * Complete the 'activityNotifications' function below.
 *
 * The function is expected to return an INTEGER.
 * The function accepts following parameters:
 *  1. INTEGER_ARRAY expenditure
 *  2. INTEGER d
 */

// Complete the activityNotifications function below.
func activityNotifications(expenditure []int, days int) int {
	nCount := 0
	mWindow := TrailingDays{}

	for exp := range expenditure {
		if exp < days {
			mWindow = mWindow.Insert(expenditure[exp])

			continue
		}

		if cutOff := mWindow.MedianTimesTwo(); expenditure[exp] >= cutOff {
			nCount++
		}

		mWindow = mWindow.Insert(expenditure[exp])
		mWindow = mWindow.Delete(expenditure[exp-days])
	}

	return nCount
}

// TrailingDays assumes elements are sorted.
type TrailingDays []int

// Insert into the slice, maintaining the sort order.
func (t TrailingDays) Insert(value int) TrailingDays {
	i := sort.SearchInts(t, value)
	n := append(t, 0)
	copy(n[i+1:], n[i:])
	n[i] = value

	return n
}

// Delete from the slice, maintaining the sort order.
func (t TrailingDays) Delete(value int) TrailingDays {
	i := sort.SearchInts(t, value)
	n := append(t[:i], t[i+1:]...)

	return n
}

// Median * 2 of the slice.
func (t TrailingDays) MedianTimesTwo() int {
	if len(t)%2 == 1 {
		return t[len(t)/2] * 2
	}

	return t[len(t)/2] + t[len(t)/2-1]
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	firstMultipleInput := strings.Split(strings.TrimSpace(readLine(reader)), " ")

	numDaysOfData, err := strconv.Atoi(firstMultipleInput[0])
	checkError(err)

	dTemp, err := strconv.Atoi(firstMultipleInput[1])
	checkError(err)

	trailingDays := dTemp

	expenditureTemp := strings.Split(strings.TrimSpace(readLine(reader)), " ")
	expenditure := make([]int, 0)

	for i := range numDaysOfData {
		expenditureItemTemp, err := strconv.Atoi(expenditureTemp[i])
		checkError(err)

		expenditureItem := expenditureItemTemp
		expenditure = append(expenditure, expenditureItem)
	}

	result := activityNotifications(expenditure, trailingDays)

	fmt.Fprintf(writer, "%d\n", result)

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
