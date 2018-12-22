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

// Complete the hourglassSum function below.
func hourglassSum(arr [][]int32) int32 {
	hourglassSums := []int{}

	for yIndex := 1; yIndex < len(arr)-1; yIndex++ {
		for xIndex := 1; xIndex < len(arr[yIndex])-1; xIndex++ {
			hourglassSums = append(hourglassSums, int(sumSingleHourglass(xIndex, yIndex, arr)))
		}
	}

	sort.Ints(hourglassSums)

	return int32(hourglassSums[len(hourglassSums)-1])
}

func sumSingleHourglass(x, y int, arr [][]int32) int32 {
	sum := int32(0)

	sum += arr[y-1][x-1]
	sum += arr[y-1][x]
	sum += arr[y-1][x+1]
	sum += arr[y][x]
	sum += arr[y+1][x-1]
	sum += arr[y+1][x]
	sum += arr[y+1][x+1]

	return sum
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024*1024)

	var arr [][]int32
	for i := 0; i < 6; i++ {
		arrRowTemp := strings.Split(readLine(reader), " ")

		var arrRow []int32
		for _, arrRowItem := range arrRowTemp {
			arrItemTemp, err := strconv.ParseInt(arrRowItem, 10, 64)
			checkError(err)
			arrItem := int32(arrItemTemp)
			arrRow = append(arrRow, arrItem)
		}

		if len(arrRow) != int(6) {
			panic("Bad input")
		}

		arr = append(arr, arrRow)
	}

	result := hourglassSum(arr)

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
