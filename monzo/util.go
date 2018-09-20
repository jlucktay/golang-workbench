package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

func convertURL(input string) *url.URL {
	// If they are just anchor links on the same page, disregard and return nil
	if strings.HasPrefix(input, "#") {
		return nil
	}

	prefix := ""

	// The assumption for relative URLs on the same domain is that they are all secure, hence prepending 'https'
	if strings.HasPrefix(input, "/") {
		if strings.HasPrefix(input, "//") {
			prefix = "https:"
		} else {
			prefix = fmt.Sprintf("https://%s", flagURL)
		}
	}

	urlOut, errParse := url.Parse(prefix + input)
	if errParse != nil {
		log.Fatalf("Error parsing '%s': %v\n", prefix+input, errParse)
	}

	// Make non-absolute URLs into absolute
	if !urlOut.IsAbs() {
		split := strings.SplitN(urlOut.String(), "/", 2)
		urlOut.Scheme = "https"

		if len(split) < 2 {
			urlOut.Host = flagURL
			urlOut.Path = "/" + split[0]
		} else {
			urlOut.Host = split[0]
			urlOut.Path = "/" + split[1]
		}
	}

	return urlOut
}
