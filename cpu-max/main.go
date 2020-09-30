// Thank you:
// https://stackoverflow.com/questions/41079492/golang-code-to-increase-cpu-usage
package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/pflag"
)

const (
	defaultDuration     = 5 * time.Minute
	progressGranularity = 100
)

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

	bar := progressbar.NewOptions(progressGranularity,
		progressbar.OptionSetDescription(fmt.Sprintf("Maxing %d CPU(s)...", runtime.NumCPU())),
		progressbar.OptionSetItsString("core(s) cooked"),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			BarStart: "[", BarEnd: "]", Saucer: "🔥", SaucerHead: "😡", SaucerPadding: "  ",
		}))

	for i := 0; i < progressGranularity; i++ {
		if err := bar.Add(1); err != nil {
			fmt.Fprintf(os.Stderr, "error adding to progress bar: %v", err)
		}

		time.Sleep(*dur / progressGranularity)
	}

	fmt.Println()
	close(done)
}
