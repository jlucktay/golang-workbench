package main

import (
	"flag"
	"fmt"
	"os"

	"./crawler"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [opts] target_url\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "    target_url should be a URL, the root of the crawl\n")
		flag.PrintDefaults()
	}

	outputFile := flag.String("output_file", "-", "If given, output to this file rather than standard out")
	verbose := flag.Bool("verbose", false, "increase output verbosity")
	threads := flag.Uint("number_of_threads", 5, "Number of goroutines to use when crawling")

	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		flag.Usage()
		return
	}

	var log crawler.Log = crawler.Silent()
	if *verbose {
		log = crawler.Verbose()
	}

	c := crawler.Crawler{*threads, log}
	if _, err := c.Crawl(args[0], *outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}
}
