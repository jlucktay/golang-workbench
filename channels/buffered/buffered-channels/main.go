package main

import (
	"fmt"
	"time"
)

func write(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Printf("Successfully wrote '%d' to ch.\n", i)
	}

	close(ch)
}

func main() {
	ch := make(chan int, 2)
	go write(ch)
	time.Sleep(2 * time.Second)

	for v := range ch {
		fmt.Printf("Read value '%d' from ch.\n", v)
		time.Sleep(2 * time.Second)
	}
}
