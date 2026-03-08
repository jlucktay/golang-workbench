package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Complete the freqQuery function below.
func freqQuery(queries [][]int32) []int32 {
	dataStructure := make(map[int32]int32, 0)
	op3Results := make([]int32, 0)

	for _, query := range queries {
		op := query[0]
		dataElement := query[1]

		switch op {
		case 1:
			dataStructure[dataElement]++

		case 2:
			if _, ok := dataStructure[dataElement]; ok {
				dataStructure[dataElement]--

				if dataStructure[dataElement] == 0 {
					delete(dataStructure, dataElement)
				}
			}

		case 3:
			found := false

			for _, v := range dataStructure {
				if v == dataElement {
					found = true

					break
				}
			}

			if found {
				op3Results = append(op3Results, 1)
			} else {
				op3Results = append(op3Results, 0)
			}

		default:
			panic(fmt.Sprintf("unknown operation: %d", op))
		}
	}

	return op3Results
}

func main() {
	for _, num := range []string{"08", "11"} {
		fmt.Printf("input: '%s', output: '%s'\n", "input"+num+".txt", "output"+num+".txt")
	}

	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	qTemp, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)
	q := int32(qTemp)

	var queries [][]int32
	for i := 0; i < int(q); i++ {
		queriesRowTemp := strings.Split(strings.TrimRight(readLine(reader), " \t\r\n"), " ")

		var queriesRow []int32
		for _, queriesRowItem := range queriesRowTemp {
			queriesItemTemp, err := strconv.ParseInt(queriesRowItem, 10, 64)
			checkError(err)
			queriesItem := int32(queriesItemTemp)
			queriesRow = append(queriesRow, queriesItem)
		}

		if len(queriesRow) != 2 {
			panic("Bad input")
		}

		queries = append(queries, queriesRow)
	}

	ans := freqQuery(queries)

	for i, ansItem := range ans {
		fmt.Fprintf(writer, "%d", ansItem)

		if i != len(ans)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}

	fmt.Fprintf(writer, "\n")

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
