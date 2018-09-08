// Thank you:
// https://stackoverflow.com/questions/41079492/golang-code-to-increase-cpu-usage
package main

import (
	"runtime"
	"time"
)

func main() {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(time.Minute * 5)
	close(done)
}
