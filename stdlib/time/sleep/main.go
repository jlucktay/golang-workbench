package main

import (
	"fmt"
	"time"
)

func main() {
	t0 := time.Now()
	time.Sleep(1 * time.Second)
	t1 := time.Now()
	fmt.Printf("The call took %v to run.\n", t1.Sub(t0))

	nowPlusOneHour := time.Now().Add(1 * time.Hour)
	fmt.Printf("One hour from now: %v\n", time.Until(nowPlusOneHour))
}
