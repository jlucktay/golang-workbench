// Thank you:
// https://stackoverflow.com/questions/41079492/golang-code-to-increase-cpu-usage
package main

import (
	"runtime"
	"time"

	"github.com/spf13/pflag"
)

const defaultDuration = 5 * time.Minute

func main() {
	dur := pflag.DurationP("duration", "d", defaultDuration, "how long to max CPU(s) for")
	pflag.Parse()

	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default: //nolint:staticcheck // That's the whole point! 😅
				}
			}
		}()
	}

	time.Sleep(*dur)
	close(done)
}
